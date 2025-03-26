// Handles authentication-related routes: register (GET/POST), login (GET/POST), logout (GET).
// Includes middleware to protect routes requiring authentication.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/sessions"
	"github.com/phenriqx/notes-api/config"
	"github.com/phenriqx/notes-api/models"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Credentials struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

type UserStore interface {
	GetUserByUsername(username string) (models.User, error)
	SaveUser(username, email, password string) error
}

type SessionStore interface {
	Get(r *http.Request, name string) (*sessions.Session, error)
	Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error
}

// Middleware is a design pattern that refers to functions or components that sit between an incoming HTTP request and the final handler that processes it.
// Middleware intercepts, processes, or modifies requests and responses as they flow through your application,
// allowing you to add reusable functionality like authentication, logging, or error handling without duplicating code in every handler.
func AuthRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		
		authHeader := r.Header.Get("Authorization")
		slog.Info("AUTH HEADER RECEIVED", "header", r.Header.Get("Authorization"))
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			slog.Error("No authorization header", "method", r.Method)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Authentication required",
			})
			return
		}

		// Get token from Authorization header (Bearer <token>)
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &Credentials{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			jwtSecret, _ := config.LoadJWTSecretKey()
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			slog.Error("Error validating JWT token", "error", err)
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{
				"error": "Invalid or expired token",
			})
			return
		}

		slog.Info("User authenticated", "user_id", claims.UserID)
		ctx := context.WithValue(r.Context(), "user_id", claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RegisterHandler(users UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var registerRequest models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
			slog.Error("Error decoding JSON register response", "error", err)
			http.Error(w, "Error decoding JSON response", http.StatusInternalServerError)
			return
		}
		if len(registerRequest.Username) == 0 || len(registerRequest.Email) == 0 || len(registerRequest.Password) == 0 {
			http.Error(w, "All fields must be filled.", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost) // Hash the provided password
		if err != nil {
			slog.Error("Error generating hashed password", "error", err)
			http.Error(w, "Error generating password.", http.StatusInternalServerError)
			return
		}

		result := users.SaveUser(
			registerRequest.Username,
			registerRequest.Email,
			string(hashedPassword),
		)
		if result != nil {
			http.Error(w, result.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User created successfully",
		})
		slog.Info("User created successfully")
	}
}

func LoginHandler(users UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var loginRequest models.LoginRequest
		err := json.NewDecoder(r.Body).Decode(&loginRequest)
		if err != nil {
			slog.Error("Error decoding JSON login response", "error", err)
			http.Error(w, "Error decoding JSON response", http.StatusInternalServerError)
			return
		}
		user, err := users.GetUserByUsername(loginRequest.Username)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				http.Error(w, "User not found", http.StatusNotFound)
				return
			}
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
			slog.Error("Passwords do not match. ERROR", "error", err)
			http.Error(w, "Passwords do not match.", http.StatusBadRequest)
			return
		}

		tokenString, err := CreateToken(user.ID)
		if err != nil {
			slog.Error("Error creating JWT token", "error", err)
            http.Error(w, "Internal server error", http.StatusInternalServerError)
            return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Logged in successfully",
			"token": tokenString,
		})
		slog.Info("User logged in succesfully")
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to handle user logout
	fmt.Fprintf(w, "Logout user")
}

func CreateToken(userID uint) (string, error) {
	JWTSecretKey, err := config.LoadJWTSecretKey()
	if err != nil {
		slog.Error("Error loading JWT key to create token", "error", err)
		return "", err
	}

	claims := &Credentials{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "notes-api",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(JWTSecretKey)
	if err != nil {
		slog.Error("Error creating JWT token", "error", err)
		return "", err
	}

	return tokenString, nil
}