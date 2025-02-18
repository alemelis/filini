package db

import (
	"fmt"
	"log"
	"os"

	"github.com/alemelis/filini/models"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// DB is the global variable representing the database connection
var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	// Define the database connection string
	if err := godotenv.Load(); err != nil {
		log.Printf("Warning: .env file not found: %v", err)
	}

	// Get environment variables with fallback values
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")
	port := os.Getenv("DB_PORT")

	// Construct the database connection string
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		host, user, password, dbname, port,
	)
	var err error

	// Initialize the DB variable
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Error connecting to the database: %v", err)
	}

	// Auto migrate models (create tables if they don't exist)
	err = DB.AutoMigrate(&models.Subtitle{})
	if err != nil {
		log.Fatalf("Error automigrating models: %v", err)
	}

	fmt.Println("Database connection established!")
}

// Insert a new subtitle into the db
func InsertSubtitle(id, videoID int, text string, startTime, endTime float64) error {
	subtitle := models.Subtitle{
		ID:        id,
		VideoID:   videoID,
		Text:      text,
		StartTime: startTime,
		EndTime:   endTime,
	}

	if err := DB.Create(&subtitle).Error; err != nil {
		return err
	}

	return nil
}

func GetVideoFilePath(videoID int) (string, error) {
	var filePath string
	result := DB.Table("videos").Select("file_path").Where("id = ?", videoID).Scan(&filePath)
	if result.Error != nil {
		return "", result.Error
	}
	return filePath, nil
}
