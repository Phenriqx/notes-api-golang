// Handles authentication-related routes: register (GET/POST), login (GET/POST), logout (GET).
// Includes middleware to protect routes requiring authentication.
package handlers

import (
	"fmt"
	"net/http"
)

func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to handle user registration
	fmt.Fprint(w, "Registe user")
}

func LoginHandler(w http.ResponseWriter, r *http.Request){
	// Implement logic to handle user login
	fmt.Fprintf(w, "Login user")
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
    // Implement logic to handle user logout
	fmt.Fprintf(w, "Logout user")
}