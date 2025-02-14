package main

import (
	"github.com/alemelis/filini/db"
	"github.com/alemelis/filini/server"
	"log"
)

func main() {
	db.InitDB()
	log.Println("Starting filini...")
	server.Start()
}
