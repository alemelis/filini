package db

import (
	"fmt"
	"github.com/alemelis/filini/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

// DB is the global variable representing the database connection
var DB *gorm.DB

// InitDB initializes the database connection
func InitDB() {
	// Define the database connection string
	dsn := "host=db user=postgres password=mysecretpassword dbname=filini port=5432 sslmode=disable"
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
