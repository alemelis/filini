package server

import (
	"database/sql"
	"fmt"
    "net/http"
	_ "github.com/lib/pq"
	"log"
	"os"
    "github.com/gin-gonic/gin"
)

var DB *sql.DB

func InitDB() {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://filini:password@localhost:5432/filini?sslmode=disable"
	}

	var err error
	DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatal("Database is unreacheable:", err)
	}

	fmt.Println("Connected to PostgreSQL!")
}

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
	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "Filini is running!"})
	})

	fmt.Println("Filini is running on port", port)
	log.Fatal(router.Run(":" + port)) // Start the server

}
