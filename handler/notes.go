// Manages note-related routes: list notes (GET /notes), create note (GET/POST /notes/new, /notes), view note (GET /notes/{id}),
// edit note (GET/POST /notes/{id}/edit, /notes/{id}), delete note (POST /notes/{id}/delete).

package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/phenriqx/notes-api/models"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

type NoteStore interface {
	FindNotesByUserID(userID uint) ([]models.Notes, error)
	DeleteNotesWithID(noteID string) error
	GetNoteByID(noteID string) (models.Notes, error)
}

func GetNotesHandler(notes NoteStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		userID, ok := r.Context().Value("user_id").(uint)
		if !ok {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": "User ID not found"})
			return
		}

		notes, err := notes.FindNotesByUserID(userID)
		if err != nil {
			http.Error(w, "Error fetching notes from database.", http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(notes)
	}
}

func CreateNoteHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		var note models.Notes
		err := json.NewDecoder(r.Body).Decode(&note)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError) // InternalServerError returns status code 500
			return
		}

		title := note.Title
		content := note.Content
		if len(title) <= 0 || len(content) <= 0 {
			http.Error(w, "Title and content must be non-empty.", http.StatusBadRequest) // BadGateway returns status code 400 - Indicates a client-side error
			return
		}

		db.Create(&note)

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "Note created successfully",
		})

		slog.Info("Note created successfully")
	}
}

func GetNoteByIDHandler(notes NoteStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		id := vars["id"]
		note, err := notes.GetNoteByID(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(note)
	}
}

func EditNoteHandler(notes NoteStore, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		id := vars["id"]

		oldNote, err := notes.GetNoteByID(id)
		if err != nil {
			http.Error(w, "Error fetching note from database.", http.StatusNotFound)
			slog.Error("Error fetching note from database", "error", err)
            return
		}

		var updatedNote models.Notes
		if err := json.NewDecoder(r.Body).Decode(&updatedNote); err != nil {
			http.Error(w, "Error decoding note.", http.StatusInternalServerError)
			slog.Error("Error decoding note", "error", err)
			return
		}

		if len(updatedNote.Title) <= 0 || len(updatedNote.Content) <= 0 {
			http.Error(w, "Title and content fields must not be empty.", http.StatusBadRequest)
			return
		}
		oldNote.Title = updatedNote.Title
		oldNote.Content = updatedNote.Content		

		db.Save(&oldNote)		
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"Message": "Note updated successfully",
		})
		
		slog.Info("Note updated successfully")
	}
}

func DeleteNoteHandler(notes NoteStore) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		id := vars["id"]
		if err := notes.DeleteNotesWithID(id); err != nil {
			http.Error(w, "Error deleting note with this specific ID.", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{
			"Message": "Note deleted successfully",
		})
	}
}