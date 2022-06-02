package main

import (
	"github.com/sMARCHz/microservices-go/hexagonal-architecture/app"
	"github.com/sMARCHz/microservices-go/hexagonal-architecture/logger"
)

func main() {
	logger.Info("Starting the application...")
	app.Start()
}
