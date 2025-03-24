package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/phenriqx/notes-api/config"
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

    gormStore := &database.GormStore{DB: db}
    sessionStore := &config.CookieSessionStore{Store: config.Sessions}

	fmt.Println("Initializing routers...")
	myRouter := mux.NewRouter()

	myRouter.Handle("/notes", handler.AuthRequiredMiddleware(sessionStore, handler.GetNotesHandler(gormStore))).Methods("GET")
	myRouter.Handle("/notes/new", handler.AuthRequiredMiddleware(sessionStore, handler.CreateNoteHandler(db))).Methods("POST")
	myRouter.Handle("/note/{id}", handler.AuthRequiredMiddleware(sessionStore, handler.GetNoteByIDHandler(db))).Methods("GET")
	myRouter.HandleFunc("/login", handler.LoginHandler(gormStore, sessionStore)).Methods("POST")
	myRouter.HandleFunc("/register", handler.RegisterHandler(gormStore)).Methods("GET", "POST")
	myRouter.HandleFunc("/logout", handler.LogoutHandler).Methods("GET", "POST")

	http.Handle("/v", myRouter)

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

/*
Why are interfaces more idiomatic?

Interfaces provide a clear and concise way to define a set of methods that a struct must implement. 
This makes it easier to work with code that expects objects of different types,
as it can be used as a placeholder for a type that has the specified methods.

Dependency injection: handlers take interfaces, not concrete types. You can swap GORM for SQLite or a mock easily

Small Interfaces: UserStore, NoteStore, SessionStore define only what’s needed, following Go’s "accept interfaces, return structs" mantra.

Testability: Mock UserStore or SessionStore for unit tests without a real DB or session store.
*/