package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/phenriqx/notes-api/handlers"
	"github.com/phenriqx/notes-api/database"
	"github.com/phenriqx/notes-api/config"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
)

func main(){
	db, err := database.Connect()
	if err != nil {
		fmt.Printf("Connection error: %v\n", err)
		return
	}

	if err := godotenv.Load(); err != nil {
		fmt.Printf("error loading godotenv: %v\n", err)
		return
	}

	config.Sessions = sessions.NewCookieStore([]byte(os.Getenv("STORE_SECRET_KEY")))
	port := os.Getenv("PORT")

	fmt.Println("Initializing routers...")
	myRouter := mux.NewRouter()

	myRouter.Handle("/notes", handlers.AuthRequiredMiddleware(handlers.GetNotesHandler(db))).Methods("GET")
	myRouter.Handle("/notes/new", handlers.AuthRequiredMiddleware(handlers.CreateNoteHandler(db))).Methods("POST")
	myRouter.Handle("/note/{id}", handlers.AuthRequiredMiddleware(handlers.GetNoteByIDHandler(db))).Methods("GET")
	myRouter.HandleFunc("/login", handlers.LoginHandler(db)).Methods("POST")
	myRouter.HandleFunc("/register", handlers.RegisterHandler(db)).Methods("GET", "POST")
	myRouter.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET", "POST")

	http.Handle("/", myRouter)

	fmt.Println("Serving on port: ", port)
	http.ListenAndServe(port, myRouter)
}