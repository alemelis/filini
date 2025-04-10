package server

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/alemelis/filini/db"
	"github.com/alemelis/filini/models"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"gorm.io/gorm"
)

func Init() {
	fmt.Println("Server package initialised")
}

func Start() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080" // Default port
	}

	r := gin.Default()
	r.Use(cors.Default())

	// Define routes
	// r.GET("/", HandleRoot)
	r.GET("/subtitles/search", HandleSearchSubtitles)
	r.Static("/storage/webm", "./storage/webm")

	for _, route := range r.Routes() {
		fmt.Println(route.Method, route.Path)
	}

	fmt.Println("Filini is running on port", port)
	log.Fatal(r.Run(":" + port)) // Start the server
}

// func HandleRoot(c *gin.Context) {
// 	c.JSON(http.StatusOK, gin.H{
// 		"message": "Filini API server is running",
// 		"endpoints": []string{
// 			"/subtitles/search?q={query}",
// 			"/storage/webm/{filename}",
// 		},
// 	})
// }

func HandleSearchSubtitles(c *gin.Context) {
	query := c.DefaultQuery("q", "")
	series := c.DefaultQuery("s", "")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	quotes, err := db.SearchSubtitles(query, series)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	if len(quotes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No subtitles found matching your query"})
	} else {
		var webms []models.Webm
		for _, quote := range quotes {
			var subtitle models.Subtitle
			if err := db.DB.First(&subtitle, quote.ID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Subtitle not found"})
				return
			}

			var video models.Video
			if err := db.DB.First(&video, subtitle.VideoID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
				return
			}

			var webm models.Webm
			if err := db.DB.First(&webm, quote.ID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Webm not found"})
				return
			}

			webms = append(webms, models.Webm{
				Model:      gorm.Model{},
				ID:         webm.ID,
				VideoId:    webm.VideoId,
				SubtitleId: webm.SubtitleId,
				FilePath:   webm.FilePath,
			})
		}
		c.JSON(http.StatusOK, webms)
	}
}
