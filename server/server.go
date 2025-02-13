package server

import (
	"fmt"
	"github.com/alemelis/filini/db"
	"github.com/alemelis/filini/models"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
	"os"
	"strings"
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
	router.GET("/subtitles/search", HandleSearchSubtitles)
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Filini is running!"})
	})

	for _, route := range router.Routes() {
	    fmt.Println(route.Method, route.Path)
	}

	fmt.Println("Filini is running on port", port)
	log.Fatal(router.Run(":" + port)) // Start the server

}

func CreateSubtitle(c *gin.Context) {
	var subtitle models.Subtitle

	// Bind the JSON data to the struct
	if err := c.ShouldBindJSON(&subtitle); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		return
	}

	// Add subtitle to the database
	err := db.InsertSubtitle(subtitle.ID, subtitle.VideoID, subtitle.Text, subtitle.StartTime, subtitle.EndTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add subtitle"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subtitle added successfully"})
}

func HandleSearchSubtitles(c *gin.Context) {
    query := c.DefaultQuery("q", "")
    if query == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
        return
    }

    var sampleSubtitles = []models.Subtitle{
        {ID: 1, VideoID: 1, Text: "This is a test subtitle.", StartTime: 0, EndTime: 5},
        {ID: 2, VideoID: 1, Text: "Another subtitle for testing.", StartTime: 5, EndTime: 10},
        {ID: 3, VideoID: 2, Text: "Some more subtitles to search.", StartTime: 0, EndTime: 5},
    }

    var results []models.Subtitle
    for _, subtitle := range sampleSubtitles {
        if strings.Contains(subtitle.Text, query) {
            results = append(results, subtitle)
        }
    }

    if len(results) == 0 {
        c.JSON(http.StatusNotFound, gin.H{"message": "No subtitles found matching your query"})
    } else {
        c.JSON(http.StatusOK, results)
    }
}
