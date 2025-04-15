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

	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"https://tuotubo.natomo.xyz", "http://localhost:3000"}
	config.AllowMethods = []string{"GET"} //, "POST", "PUT", "PATCH", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Length", "Content-Type", "Authorization"}
	config.ExposeHeaders = []string{"Content-Length", "Content-Range", "Accept-Ranges"}

	r := gin.Default()
	r.Use(cors.New(config))

	// Define routes
	r.GET("/subtitles/search", HandleSearchSubtitles)
	r.Static("/storage/webm", "./storage/webm")
	r.Static("/storage/poster", "./storage/poster")

	for _, route := range r.Routes() {
		fmt.Println(route.Method, route.Path)
	}

	fmt.Println("Filini is running on port", port)
	log.Fatal(r.Run(":" + port)) // Start the server
}

func HandleSearchSubtitles(c *gin.Context) {
	query := c.DefaultQuery("q", "")
	series := c.DefaultQuery("s", "")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	if series == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 's' is required"})
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
		var clips []models.Clip
		for _, quote := range quotes {
			var subtitle models.Subtitle
			if err := db.DB.First(&subtitle, quote.ID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Subtitle not found"})
				return
			}

			var webm models.Webm
			if err := db.DB.First(&webm, quote.ID).Error; err != nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Webm not found"})
				return
			}

			clips = append(clips, models.Clip{
				Model:    gorm.Model{},
				Quote:    subtitle.Text,
				WebmPath: webm.FilePath,
			})
		}
		c.JSON(http.StatusOK, clips)
	}
}
