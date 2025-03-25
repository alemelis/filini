package server

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"

	"github.com/alemelis/filini/db"
	"github.com/alemelis/filini/models"
	"github.com/alemelis/filini/utils"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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

	r := gin.Default()
	r.Use(cors.Default())

	// Define routes
	r.POST("/upload", HandleUpload)
	r.POST("/generate_gif", HandleGenerateGIF)
	r.GET("/subtitles/search", HandleSearchSubtitles)
	r.GET("/gif/:subtitle_id", HandleGenerateGIF)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Filini is running!"})
	})
	r.Static("/storage/gifs", "./storage/gifs")

	for _, route := range r.Routes() {
		fmt.Println(route.Method, route.Path)
	}

	fmt.Println("Filini is running on port", port)
	log.Fatal(r.Run(":" + port)) // Start the server
}

func HandleUpload(c *gin.Context) {
	id := uuid.New().ID()

	title := c.PostForm("title")
	if title == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get title"})
		return
	}

	video_path := c.PostForm("video_path")
	if video_path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get video_path"})
		return
	}

	srt_path := c.PostForm("srt_path")
	if srt_path == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get srt_path"})
		return
	}

	// video
	var err error
	err = db.InsertVideo(id, title, video_path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload video"})
		return
	}

	// subtitles
	HandleUploadSubtitles(c, id, title, srt_path)

	c.JSON(http.StatusOK, gin.H{"message": "Upload successful"})
}

func HandleUploadSubtitles(c *gin.Context, videoID uint32, title string, srt_path string) {
	file, err := os.Open(srt_path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open SRT file"})
		return
	}
	defer file.Close()

	// Parse SRT file
	subtitles, err := utils.ParseSRT(file, videoID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse SRT file"})
		return
	}

	// Insert subtitles into DB
	if err := db.InsertSubtitles(subtitles); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store subtitles"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Subtitles uploaded successfully"})
}

// func CreateVideo(c *gin.Context) {
// 	id, err := strconv.Atoi(c.PostForm("id"))
// 	if err != nil {
// 		id = uuid.New().String()
// 	}

// 	title := c.PostForm("title")
// 	if title == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get title"})
// 		return
// 	}

// 	file_path := c.PostForm("file_path")
// 	if file_path == "" {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file_path"})
// 		return
// 	}

// 	// Add video to the database
// 	err = db.InsertVideo(id, title, file_path)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add video"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Video added successfully"})
// }

// func CreateSubtitle(c *gin.Context) {
// 	var subtitle models.Subtitle

// 	// Bind the JSON data to the struct
// 	if err := c.ShouldBindJSON(&subtitle); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	// Add subtitle to the database
// 	err := db.InsertSubtitle(subtitle.VideoID, subtitle.Text, subtitle.StartTime, subtitle.EndTime)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add subtitle"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Subtitle added successfully"})
// }

// func CreateSubtitles(c *gin.Context) {
// 	var subtitles []models.Subtitle

// 	// Bind JSON array of subtitles
// 	if err := c.ShouldBindJSON(&subtitles); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
// 		return
// 	}

// 	// Insert all subtitles into the database
// 	if err := db.InsertSubtitles(subtitles); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add subtitles"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Subtitles added successfully"})
// }

// func HandleUploadSubtitles(c *gin.Context) {
// 	videoID, err := strconv.Atoi(c.PostForm("video_id"))
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid or missing video_id"})
// 		return
// 	}

// 	file, err := c.FormFile("file")
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to get file"})
// 		return
// 	}

// 	src, err := file.Open()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
// 		return
// 	}
// 	defer src.Close()

// 	// Parse SRT file
// 	subtitles, err := utils.ParseSRT(src, videoID)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse SRT file"})
// 		return
// 	}

// 	// Insert subtitles into DB
// 	if err := db.InsertSubtitles(subtitles); err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to store subtitles"})
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Subtitles uploaded successfully"})
// }

func HandleSearchSubtitles(c *gin.Context) {
	query := c.DefaultQuery("q", "")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Query parameter 'q' is required"})
		return
	}

	quotes, err := db.SearchSubtitles(query)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database query failed"})
		return
	}

	if len(quotes) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No subtitles found matching your query"})
	} else {
		var gifs []models.Gif
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

			// gifPath := generateGIF(video.ID, subtitle.ID)
			// gifs = append(gifs, models.Gif{
			// 	Quote:    subtitle.Text,
			// 	FilePath: gifPath,
			// })
		}
		c.JSON(http.StatusOK, gifs)
	}
}

func generateGIF(videoID, subtitleID uint32) string {
	// Define GIF output path
	gifPath := fmt.Sprintf("storage/gifs/td-%d-%d.gif", videoID, subtitleID)

	if _, err := os.Stat(gifPath); !os.IsNotExist(err) {
		return gifPath
	}

	var subtitle models.Subtitle
	db.DB.First(&subtitle, subtitleID)

	// Fetch the associated video
	var video models.Video
	db.DB.First(&video, subtitle.VideoID)

	// Run ffmpeg to generate GIF
	cmd := exec.Command("ffmpeg", "-i", video.FilePath, "-ss", fmt.Sprintf("%f", subtitle.StartTime), "-to", fmt.Sprintf("%f", subtitle.EndTime+1.0), "-vf", "fps=10,scale=320:-1", "-err_detect", "ignore_err", "-fflags", "+genpts+discardcorrupt", "-threads", "1", "-y", gifPath)

	// Create pipes for stdout and stderr
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("Error creating stdout pipe: %v\n", err)
		return ""
	}
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Printf("Error creating stderr pipe: %v\n", err)
		return ""
	}

	// Start the command
	if err := cmd.Start(); err != nil {
		log.Printf("Error starting command: %v\n", err)
		return ""
	}

	// Read from pipes in goroutines
	go func() {
		scanner := bufio.NewScanner(stdoutPipe)
		for scanner.Scan() {
			log.Printf("FFmpeg stdout: %s\n", scanner.Text())
		}
	}()

	go func() {
		scanner := bufio.NewScanner(stderrPipe)
		for scanner.Scan() {
			log.Printf("FFmpeg stderr: %s\n", scanner.Text())
		}
	}()

	// Wait for the command to complete
	if err := cmd.Wait(); err != nil {
		log.Printf("FFmpeg error: %v\n", err)
		return ""
	}

	return gifPath
}

func HandleGenerateGIF(c *gin.Context) {
	subtitleID := c.Param("subtitle_id")

	// Fetch subtitle details
	var subtitle models.Subtitle
	if err := db.DB.First(&subtitle, subtitleID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Subtitle not found"})
		return
	}

	// Fetch the associated video
	var video models.Video
	if err := db.DB.First(&video, subtitle.VideoID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
		return
	}

	// Define GIF output path
	gifPath := fmt.Sprintf("storage/gifs/td-%d-%d.gif", video.ID, subtitle.ID)

	// Run ffmpeg to generate GIF
	cmd := exec.Command("ffmpeg", "-i", video.FilePath, "-ss", fmt.Sprintf("%f", subtitle.StartTime), "-to", fmt.Sprintf("%f", subtitle.EndTime+1.0), "-vf", "fps=10,scale=320:-1", "-y", gifPath)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate GIF", "details": stderr.String()})
		return
	}

	// Return the generated GIF path
	c.JSON(http.StatusOK, gin.H{"gif_url": fmt.Sprintf("/%s", gifPath)})
}
