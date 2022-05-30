package main

import (
	"github.com/sMARCHz/microservices-go/app"
	"github.com/sMARCHz/microservices-go/logger"
)

func main() {
	logger.Info("Starting the application...")
	app.Start()
}
