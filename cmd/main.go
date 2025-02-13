package main

import (
    "log"
    "github.com/alemelis/filini/server"
    "github.com/alemelis/filini/db"
)

func main() {
    db.InitDB()
    log.Println("Starting filini...")
    server.Start()
}

