package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/phenriqx/notes-api/database"
	"github.com/phenriqx/notes-api/handler"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	"github.com/jub0bs/cors"
)

func main() {
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

	myRouter.Handle("/notes", handler.AuthRequiredMiddleware(handler.GetNotesHandler(db))).Methods("GET")
	myRouter.Handle("/notes/new", handler.AuthRequiredMiddleware(handler.CreateNoteHandler(db))).Methods("POST")
	myRouter.Handle("/note/{id}", handler.AuthRequiredMiddleware(handler.GetNoteByIDHandler(db))).Methods("GET")
	myRouter.HandleFunc("/login", handler.LoginHandler(db)).Methods("POST")
	myRouter.HandleFunc("/register", handler.RegisterHandler(db)).Methods("GET", "POST")
	myRouter.HandleFunc("/logout", handler.LogoutHandler).Methods("GET", "POST")

	http.Handle("/", myRouter)

	corsMw, err := cors.NewMiddleware(cors.Config{
		Origins:        []string{"http://localhost:8080"},
		Methods:        []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodDelete, http.MethodPut},
		RequestHeaders: []string{"Content-Type"},
	})
	if err != nil {
		log.Fatal(err)
	}

	corsMw.SetDebug(true)

	fmt.Println("Serving on port: ", port)
	http.ListenAndServe(port, corsMw.Wrap(myRouter))
}
