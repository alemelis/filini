package main

import (
    "log"
    "github.com/alemelis/filini/server"
)

func main() {
    server.InitDB()
    log.Println("Starting filini...")
    server.Start()
}

