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
	// err = DB.AutoMigrate(&models.Subtitle{})
	// if err != nil {
	// 	log.Fatalf("Error automigrating models: %v", err)
	// }

	// err = DB.AutoMigrate(&models.Video{})
	// if err != nil {
	// 	log.Fatalf("Error automigrating models: %v", err)
	// }

	// err = DB.AutoMigrate(&models.Webm{})
	// if err != nil {
	// 	log.Fatalf("Error automigrating models: %v", err)
	// }

	fmt.Println("Database connection established!")
}

func SearchSubtitles(query string, series string) ([]models.Subtitle, error) {
	var video_ids []int
	DB.Table("videos").Select("id").Where("series = ?", series).Find(&video_ids)

	var results []models.Subtitle
	for _, video_id := range video_ids {
		var video_results []models.Subtitle
		err := DB.Table("subtitles").Where("video_id = ?", video_id).Where("text ILIKE ?", "%"+query+"%").Find(&video_results).Error
		if err != nil {
			return nil, err
		}
		for _, video_result := range video_results {
			results = append(results, video_result)
		}
	}
	return results, nil
}
