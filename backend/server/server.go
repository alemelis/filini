package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alemelis/filini/db"
	"github.com/alemelis/filini/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
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
	router.Use(cors.Default())

	// Define routes
	router.POST("/subtitles", CreateSubtitle)
	router.GET("/subtitles/search", HandleSearchSubtitles)
	router.POST("/generate_gif", GenerateGIFHandler)
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

type GenerateGIFRequest struct {
	VideoID   int     `json:"video_id"`
	StartTime float64 `json:"start_time"`
	EndTime   float64 `json:"end_time"`
}

func GenerateGIFHandler(c *gin.Context) {
	var req GenerateGIFRequest

	// Parse JSON request
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Validate input
	if req.EndTime <= req.StartTime {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid time range"})
		return
	}

	// Get the video file path from the database
	videoPath, err := db.GetVideoFilePath(req.VideoID) // Fixed: Use db package to get video path
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Generate the GIF file name
	gifFilename := fmt.Sprintf("gif_%d_%d.gif", req.VideoID, time.Now().Unix())
	gifPath := filepath.Join("/tmp", gifFilename)

	// Call ffmpeg to generate the GIF
	err = generateGIF(videoPath, req.StartTime, req.EndTime, gifPath)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate GIF"})
		return
	}

	// Respond with the GIF URL
	c.JSON(http.StatusOK, gin.H{"gif_url": "/tmp/" + gifFilename})
}

func generateGIF(videoPath string, startTime, endTime float64, outputPath string) error {
	cmd := exec.Command("ffmpeg",
		"-i", videoPath,
		"-ss", fmt.Sprintf("%.2f", startTime),
		"-to", fmt.Sprintf("%.2f", endTime),
		"-vf", "fps=10,scale=320:-1",
		"-y", outputPath,
	)

	return cmd.Run()
}
