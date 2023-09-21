package main

import (
	"log"
	"transactions_app/config"
	"transactions_app/router"
	"transactions_app/server"
)

func main() {

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	router := router.InitializeRouter()
	server := server.NewServer(":"+config.Port, router)
	server.Start()
}
