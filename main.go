package main

import (
	"log"

	"github.com/fiskaly/coding-challenges/signing-service-challenge/api"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/persistence"
	"github.com/fiskaly/coding-challenges/signing-service-challenge/service"
)

const (
	ListenAddress = ":8080"
)

func main() {
	repository := persistence.NewInMemoryRepository()
	deviceService := service.NewDeviceService(repository)
	transactionService := service.NewTransactionService(repository)
	server := api.NewServer(ListenAddress, deviceService, transactionService)

	if err := server.Run(); err != nil {
		log.Fatal("Could not start server on ", ListenAddress)
	}
}
