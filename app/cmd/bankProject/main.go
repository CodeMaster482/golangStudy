package main

import (
	"log"

	"go-web-bmstu/internal/api"
)

func main() {
	log.Println("App running")

	api.StartServer()

	log.Println("localhost:8080")

	log.Println("App terminated")
}
