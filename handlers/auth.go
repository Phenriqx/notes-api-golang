// Handles authentication-related routes: register (GET/POST), login (GET/POST), logout (GET).
// Includes middleware to protect routes requiring authentication.
package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/phenriqx/notes-api/models"
	"github.com/phenriqx/notes-api/config"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Middleware is a design pattern that refers to functions or components that sit between an incoming HTTP request and the final handler that processes it.
// Middleware intercepts, processes, or modifies requests and responses as they flow through your application, 
// allowing you to add reusable functionality like authentication, logging, or error handling without duplicating code in every handler.
func AuthRequiredMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := config.Sessions.Get(r, "auth-session")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if session.Values["user_id"] == nil {
			http.Error(w, fmt.Sprintf("Error: Authentication Required.\nStatus Code: %d", http.StatusUnauthorized), http.StatusUnauthorized)
            return
		}

		next.ServeHTTP(w, r)
	})
}

func RegisterHandler(db *gorm.DB) http.HandlerFunc {
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

		user := models.User{
			Username: registerRequest.Username,
			Email:    registerRequest.Email,
			Password: string(hashedPassword),
		}
		if err := db.Where("username = ? OR email = ?", user.Username, user.Email).
			First(&user).Error; err == nil {
			http.Error(w, "This user already exists in the database", http.StatusConflict) // Returns status code 409 indicating that the request could not be completed
			return
		}

		if result := db.Create(&user); result.Error != nil {
			http.Error(w, result.Error.Error(), http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "User created successfully",
		})
		log.Println("User created successfully!")
	}
}

func LoginHandler(db *gorm.DB) http.HandlerFunc{
	return func (w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var loginRequest models.LoginRequest
		err := json.NewDecoder(r.Body).Decode(&loginRequest)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		var user models.User
		result := db.Where("username = ?", loginRequest.Username).First(&user)
		err = result.Error
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "User not found.", http.StatusNotFound)
			return
		}
		
		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password)); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return	
		}

		session, err := config.Sessions.Get(r, "auth-session") // Using sessions to keep  track of logged in users => The browser gets a cookie named auth-session
		if err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
		session.Values["userId"] = user.ID
		session.Save(r, w)
		session.Options.Secure = false

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
