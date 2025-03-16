// Manages note-related routes: list notes (GET /notes), create note (GET/POST /notes/new, /notes), view note (GET /notes/{id}),
// edit note (GET/POST /notes/{id}/edit, /notes/{id}), delete note (POST /notes/{id}/delete).

package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetNotesHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to retrieve notes from the database and display them in a list
	fmt.Fprint(w, "Getting notes from database")
}

func CreateNoteHandler(w http.ResponseWriter, r *http.Request) {
	// Implement logic to create a new note
	fmt.Printf("Craete note")
}

func GetNoteByIDHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fmt.Fprintf(w, "User ID: %s\n", id)
}
