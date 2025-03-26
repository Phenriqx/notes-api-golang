package main

import (
	"fmt"
	"log/slog"
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
	// Loading and importing environment variables + connecting database
	db, err := database.Connect()
	if err != nil {
		fmt.Printf("Connection error: %v\n", err)
		return
	}
	if err := godotenv.Load(); err != nil {
		fmt.Printf("error loading godotenv: %v\n", err)
		return
	}
	PORT := os.Getenv("PORT")
	LOG_FILE := os.Getenv("LOG_FILE")

	// Interface settings
	gormStore := &database.GormStore{DB: db}
	sessionStore := &config.CookieSessionStore{Store: config.Sessions}

	// Logging configuration
	handlerOptions := &slog.HandlerOptions{
		Level:      slog.LevelDebug,
	}
	file, err := os.OpenFile(LOG_FILE, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		slog.Error("Error opening log file", "error", err)
		return
	}
	defer file.Close()
	logger := slog.New(slog.NewJSONHandler(file, handlerOptions))
	slog.SetDefault(logger)

	myRouter := mux.NewRouter()
	// Note routes
	myRouter.Handle("/notes", handler.AuthRequiredMiddleware(sessionStore, handler.GetNotesHandler(gormStore))).Methods("GET")
	myRouter.Handle("/notes/new", handler.AuthRequiredMiddleware(sessionStore, handler.CreateNoteHandler(db))).Methods("POST")
	myRouter.Handle("/note/{id}", handler.AuthRequiredMiddleware(sessionStore, handler.GetNoteByIDHandler(gormStore))).Methods("GET")
	myRouter.Handle("/note/{id}/delete", handler.AuthRequiredMiddleware(sessionStore, handler.DeleteNoteHandler(gormStore))).Methods("POST")
	myRouter.Handle("/note/{id}/update", handler.AuthRequiredMiddleware(sessionStore, handler.EditNoteHandler(gormStore, db))).Methods("POST")

	// Authentication routes
	myRouter.HandleFunc("/login", handler.LoginHandler(gormStore, sessionStore)).Methods("POST")
	myRouter.HandleFunc("/register", handler.RegisterHandler(gormStore)).Methods("GET", "POST")
	myRouter.HandleFunc("/logout", handler.LogoutHandler).Methods("GET", "POST")

	http.Handle("/", myRouter)

	corsMw, err := cors.NewMiddleware(cors.Config{
		Origins:        []string{"http://localhost:8080"},
		Methods:        []string{http.MethodGet, http.MethodPost, http.MethodOptions, http.MethodDelete, http.MethodPut},
		RequestHeaders: []string{"Content-Type"},
	})
	if err != nil {
		slog.Error("Error handling CORS: ", "error", err)
		return
	}

	corsMw.SetDebug(true)
	http.ListenAndServe(PORT, corsMw.Wrap(myRouter))
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
