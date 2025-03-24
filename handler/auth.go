// Handles authentication-related routes: register (GET/POST), login (GET/POST), logout (GET).
// Includes middleware to protect routes requiring authentication.
package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/phenriqx/notes-api/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserStore interface {
	GetUserByUsername(username string) (models.User, error)
	SaveUser(username, email, password string) error
}

type SessionStore interface {
	Get(r *http.Request, name string) (*sessions.Session, error)
	Save(r *http.Request,  w http.ResponseWriter, session *sessions.Session) error
}

// Middleware is a design pattern that refers to functions or components that sit between an incoming HTTP request and the final handler that processes it.
// Middleware intercepts, processes, or modifies requests and responses as they flow through your application,
// allowing you to add reusable functionality like authentication, logging, or error handling without duplicating code in every handler.
func AuthRequiredMiddleware(sessions SessionStore, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		session, err := sessions.Get(r, "auth-session")
		if err != nil {
			log.Printf("Session error: %v", err)
			http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			return
		}

		userID, ok := session.Values["user_id"]
		if !ok {
			log.Println("No user_id found in session")
			http.Error(w, `{"error":"Internal server error"}`, http.StatusInternalServerError)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user_id", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func RegisterHandler(users UserStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var registerRequest models.RegisterRequest
		if err := json.NewDecoder(r.Body).Decode(&registerRequest); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if len(registerRequest.Username) == 0 || len(registerRequest.Email) == 0 || len(registerRequest.Password) == 0 {
			http.Error(w, "All fields must be filled.", http.StatusBadRequest)
			return
		}

		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(registerRequest.Password), bcrypt.DefaultCost) // Hash the provided password
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
		log.Println("User created successfully!")
	}
}

func LoginHandler(users UserStore, sessions SessionStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var loginRequest models.LoginRequest
		err := json.NewDecoder(r.Body).Decode(&loginRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		session, err := sessions.Get(r, "auth-session")
		if err != nil {
			log.Printf("Error getting session: %v", err)
            http.Error(w, "Failed to get session.", http.StatusInternalServerError)
            return
		}

		session.Values["user_id"] = user.ID
		if err := session.Save(r, w); err != nil {
			log.Printf("Error saving session: %v", err)
			http.Error(w, "Failed to save session.", http.StatusInternalServerError)
			return
		}

		log.Println("Session saved successfully")

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Logged in successfully",
		})
		log.Println("Logged in successfully!")
	}
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to handle user logout
	fmt.Fprintf(w, "Logout user")
}
