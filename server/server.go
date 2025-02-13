package server

import (
	"fmt"
    "net/http"
	_ "github.com/lib/pq"
	"log"
	"os"
    "github.com/gin-gonic/gin"
    "github.com/alemelis/filini/db"
)

func Init() {
	fmt.Println("Server package initialised")
}

func Start() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

    router := gin.Default()

	// Define routes
    router.POST("/subtitles", CreateSubtitle)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Filini is running!"})
	})

	fmt.Println("Filini is running on port", port)
	log.Fatal(router.Run(":" + port)) // Start the server

}

func CreateSubtitle(c *gin.Context) {
    // Struct to bind the JSON request body to a Go object
    var subtitle struct {
        VideoID int `json:"video_id"`
        Text string `json:"text"`
        StartTime float64 `json:"start_time"`
        EndTime float64 `json:"end_time"`
    }

    // Bind the JSON data to the struct
    if err := c.ShouldBindJSON(&subtitle); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    // Add subtitle to the database
    err := db.InsertSubtitle(subtitle.VideoID, subtitle.Text, subtitle.StartTime, subtitle.EndTime)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add subtitle"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Subtitle added successfully"})
}
