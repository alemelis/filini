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

	err = DB.AutoMigrate(&models.Video{})
	if err != nil {
		log.Fatalf("Error automigrating models: %v", err)
	}

	err = DB.AutoMigrate(&models.Gif{})
	if err != nil {
		log.Fatalf("Error automigrating models: %v", err)
	}

	fmt.Println("Database connection established!")
}

func InsertVideo(id uint32, title string, filePath string) error {
	video := models.Video{
		Model:    gorm.Model{},
		ID:       id,
		Title:    title,
		FilePath: filePath,
	}

	if err := DB.Create(&video).Error; err != nil {
		return err
	}

	return nil
}

func InsertSubtitle(videoID uint32, text string, startTime, endTime float64) error {
	subtitle := models.Subtitle{
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

func InsertSubtitles(subtitles []models.Subtitle) error {
	if len(subtitles) == 0 {
		return nil
	}
	return DB.Create(&subtitles).Error
}

func SearchSubtitles(query string) ([]models.Subtitle, error) {
	var results []models.Subtitle
	err := DB.Table("subtitles").Where("text ILIKE ?", "%"+query+"%").Find(&results).Error
	if err != nil {
		return nil, err
	}
	return results, nil
}

func GetVideoFilePath(videoID int) (string, error) {
	var filePath string
	result := DB.Table("videos").Select("file_path").Where("id = ?", videoID).Scan(&filePath)
	if result.Error != nil {
		return "", result.Error
	}
	return filePath, nil
}
