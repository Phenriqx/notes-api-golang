package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/phenriqx/notes-api/handlers"
	"github.com/phenriqx/notes-api/db"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main(){
	db.Connect()

	if err := godotenv.Load(); err != nil {
		fmt.Printf("error loading godotenv: %v\n", err)
		return
	}
	port := os.Getenv("PORT")

	fmt.Println("Initializing routers...")
	myRouter := mux.NewRouter()

	myRouter.HandleFunc("/notes", handlers.GetNotesHandler).Methods("GET")
	myRouter.HandleFunc("/notes/new", handlers.CreateNoteHandler).Methods("POST")
	myRouter.HandleFunc("/notes/{id}", handlers.GetNoteByIDHandler).Methods("GET")
	myRouter.HandleFunc("/login", handlers.LoginHandler).Methods("POST")
	myRouter.HandleFunc("/logout", handlers.RegisterHandler).Methods("POST")
	myRouter.HandleFunc("/register", handlers.LogoutHandler).Methods("POST")

	http.Handle("/", myRouter)

	fmt.Println("Serving on port: ", port)
	http.ListenAndServe(port, myRouter)
}