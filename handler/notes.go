// Manages note-related routes: list notes (GET /notes), create note (GET/POST /notes/new, /notes), view note (GET /notes/{id}),
// edit note (GET/POST /notes/{id}/edit, /notes/{id}), delete note (POST /notes/{id}/delete).

package handler

import (
	"encoding/json"
	"net/http"

	"github.com/phenriqx/notes-api/models"
	"gorm.io/gorm"

	"github.com/gorilla/mux"
)

func GetNotesHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		var notes []models.Notes
		db.Find(&notes)
		userID, ok := r.Context().Value("user_id").(uint) // Adjust type based on models.User.ID
		if !ok {
			http.Error(w, "User ID not found in context", http.StatusInternalServerError)
			return
		}
        if err := db.Where("user_id = ?", userID).Find(&notes).Error; err != nil {
            json.NewEncoder(w).Encode(map[string]string{"error": "Failed to fetch notes"})
            w.WriteHeader(http.StatusInternalServerError)
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
	}
}

func GetNoteByIDHandler(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		vars := mux.Vars(r)
		id := vars["id"]
		var note models.Notes
		if err := db.Where("id = ?", id).First(&note).Error; err != nil {
			http.Error(w, "Error fetching note.", http.StatusNotFound)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(note)
	}
}
