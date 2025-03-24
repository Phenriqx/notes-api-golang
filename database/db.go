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

func (db *GormStore) FindNotesByUserID(userID uint) ([]models.Notes, error) {
	var notes []models.Notes
	if err := db.DB.Where("user_id = ?", userID).First(&notes).Error; err != nil {
		log.Println("Error getting notes from database: ", err)
		return nil, err
	} 
	return notes, nil
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
