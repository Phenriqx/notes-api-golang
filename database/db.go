package database

import (
	"fmt"
	"log"
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
		log.Println("Error getting user from database: ", err)
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
		log.Printf("User already exists in the database: %v", err)
		return err
	}

	if result := db.DB.Create(&newUser); result.Error != nil {
		log.Printf("Error creating user in the database: %v", result.Error)
        return result.Error
	}

	return nil
}

func (db *GormStore) FindNotesByUserID(userID uint) ([]models.Notes, error) {
	var notes []models.Notes
	if err := db.DB.Where("user_id = ?", userID).Find(&notes).Error; err != nil {
		log.Println("Error getting notes from database: ", err)
		return nil, err
	}

	return notes, nil
}

func (db *GormStore) DeleteNotesWithID(noteID string) error {
	var note models.Notes
	if err := db.DB.Where("id = ?", noteID).First(&note).Error; err != nil {
		log.Printf("Error getting note from database: %v", err)
		return err
	}

	if err := db.DB.Delete(&note).Error; err != nil {
		log.Printf("Error deleting note from database: %v", err)
	    return err
	}

	log.Printf("Note deleted from database with ID: %s", noteID)
	return nil
}

func Connect() (*gorm.DB, error) {
	if err := godotenv.Load(); err != nil {
		fmt.Printf("error loading godotenv: %v\n", err)
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
		fmt.Printf("Error connecting to the database: %v\n", err)
		return nil, err
	}
	fmt.Println("Database connection successful!")

	userExists := db.Migrator().HasTable(&models.User{})
	notesExists := db.Migrator().HasTable(&models.Notes{})
	if !notesExists || !userExists {
		fmt.Println("Table does not exist! Performing migrations...")
		db.AutoMigrate(&models.User{}, &models.Notes{})
	} else {
		fmt.Println("Tables exist! No migrations needed.")
		fmt.Println()
	}

	return db, nil
}
