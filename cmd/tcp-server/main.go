package main

import (
	"log"
	"mangahub/internal/tcp"
)

func main() {
	port := "9090"
	server := tcp.NewServer(port)

	log.Println("Starting MangaHub TCP Sync Server...")
	log.Printf("Listening on port %s...", port)

	server.Start()
}
