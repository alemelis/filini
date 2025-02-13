package main

import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
    // "github.com/alemelis/filini/server"
)

func main() {
    r := gin.Default()

    // Define routes
    r.GET("/", func(c *gin.Context) {
        c.JSON(http.StatusOK, gin.H{"message": "Filini is alive!"})
    })

    // Start server
    log.Println("Starting server on :8080...")
    if err := r.Run(":8080"); err != nil {
        log.Fatal(err)
    }
}
