package main

import (
	"github.com/sMARCHz/rest-based-microservices-go/hexagonal-architecture/app"
	"github.com/sMARCHz/rest-based-microservices-go/hexagonal-architecture/logger"
)

func main() {
	logger.Info("Starting the application...")
	app.Start()
}
