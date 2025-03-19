package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/phenriqx/notes-api/handlers"
	"github.com/phenriqx/notes-api/database"

	"github.com/gorilla/mux"
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
	port := os.Getenv("PORT")

	fmt.Println("Initializing routers...")
	myRouter := mux.NewRouter()

	myRouter.HandleFunc("/notes", handlers.GetNotesHandler(db)).Methods("GET")
	myRouter.HandleFunc("/notes/new", handlers.CreateNoteHandler(db)).Methods("POST")
	myRouter.HandleFunc("/note/{id}", handlers.GetNoteByIDHandler).Methods("GET")
	myRouter.HandleFunc("/login", handlers.LoginHandler).Methods("GET", "POST")
	myRouter.HandleFunc("/register", handlers.RegisterHandler(db)).Methods("GET", "POST")
	myRouter.HandleFunc("/logout", handlers.LogoutHandler).Methods("GET", "POST")

	http.Handle("/", myRouter)

	fmt.Println("Serving on port: ", port)
	http.ListenAndServe(port, myRouter)
}