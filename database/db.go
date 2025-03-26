package database

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/phenriqx/notes-api/models"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type GormStore struct {
	DB *gorm.DB
}

func (db *GormStore) GetUserByUsername(username string) (models.User, error) {
	var user models.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("User not found in database.")
			return models.User{}, err
		}
		slog.Error("Error getting user from database", "error", err)
		return models.User{}, err
	}

	return user, nil
}

func (db *GormStore) SaveUser(username, email, password string) error {
	newUser := models.User{
		Username: username,
		Email:    email,
		Password: password,
	}

	if err := db.DB.Where("username = ? OR email = ?", username, email).First(&newUser).Error; err == nil {
		slog.Error("User already exists in the database:",  "error", err)
		return err
	}

	if result := db.DB.Create(&newUser); result.Error != nil {
		slog.Error("Error creating user in the database: ", "error", result)
		return result.Error
	}

	return nil
}

func (db *GormStore) FindNotesByUserID(userID uint) ([]models.Notes, error) {
	var notes []models.Notes
	if err := db.DB.Where("user_id = ?", userID).Find(&notes).Error; err != nil {
		slog.Error("Error getting notes from database: ", "error", err)
		return nil, err
	}

	return notes, nil
}

func (db *GormStore) DeleteNotesWithID(noteID string) error {
	var note models.Notes
	if err := db.DB.Where("id = ?", noteID).First(&note).Error; err != nil {
		slog.Error("Error getting note from database: ", "error", err)
		return err
	}

	if err := db.DB.Delete(&note).Error; err != nil {
		slog.Error("Error deleting note from database: ", "error", err)
		return err
	}

	slog.Info("Note deleted succesfully.")
	return nil
}

func (db *GormStore) GetNoteByID(noteID string) (models.Notes, error) {
	var note models.Notes
	if err := db.DB.Where("id = ?", noteID).First(&note).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			slog.Error("Record not found.")
			return models.Notes{}, err
		}
		slog.Error("Error getting note from database: ", "error", err)
		return models.Notes{}, err
	}

	return note, nil
}

func Connect() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		slog.Error("Error loading godotenv: ", "error", err)
		return nil, err
	}

	host := os.Getenv("HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("PASSWORD")
	db_name := os.Getenv("DB_NAME")
	db_port := os.Getenv("DB_PORT")

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, db_name, db_port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		slog.Error("Error connecting to the database", "error", err)
		return nil, err
	}

	userExists := db.Migrator().HasTable(&models.User{})
	notesExists := db.Migrator().HasTable(&models.Notes{})
	if !notesExists || !userExists {
		slog.Error("Table does not exist! Performing migrations...", "error", err)
		db.AutoMigrate(&models.User{}, &models.Notes{})
	} else {
		fmt.Println("Tables exist! No migrations needed.")
	}

	return db, nil
}