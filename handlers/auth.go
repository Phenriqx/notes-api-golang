// Handles authentication-related routes: register (GET/POST), login (GET/POST), logout (GET).
// Includes middleware to protect routes requiring authentication.
package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/phenriqx/notes-api/models"

	// "github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

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

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to handle user login
	fmt.Fprintf(w, "Login user")
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to handle user logout
	fmt.Fprintf(w, "Logout user")
}
