package main

import (
	"log"
	"transactions_app/config"
	"transactions_app/database"
	"transactions_app/handlers"
	"transactions_app/repository"
	"transactions_app/router"
	"transactions_app/server"
	"transactions_app/usecase"
)

func main() {

	config, err := config.LoadConfig(".")
	if err != nil {
		log.Fatalf("could not load config: %v", err)
	}

	instance := database.ConnectSingleton(config.PostgresConn)
	database.Migrate(instance)

	transactionRepo := repository.NewTransactionRepositoryPsql(instance)

	transactionUsecase := usecase.NewTransactionUsecase(transactionRepo)

	transactionHandler := handlers.NewTransactionHandler(transactionUsecase, config)

	router := router.InitializeRouter(transactionHandler)
	server := server.NewServer(":"+config.Port, router)
	server.Start()
}
