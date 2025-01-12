package main

import (
	"context"
	"cynxhostagent/internal/app"
	"cynxhostagent/internal/controller"
	"flag"
	"log"

	"github.com/joho/godotenv"
)

func main() {
	log.Println("Starting CynxHost")
	ctx, cancel := context.WithCancel(context.Background())
	defer func() {
		cancel()
	}()

	log.Println("Parsing .env")
	envFile := flag.String("env", ".env", ".env")
	flag.Parse()

	// Load the specified .env file
	log.Println("Loading .env")
	err := godotenv.Load(*envFile)
	if err != nil {
		panic(err)
	}

	log.Println("Initializing App")
	app, err := app.NewApp(ctx, "config.json")
	if err != nil {
		panic(err)
	}

	logger := app.Dependencies.Logger

	logger.Infoln("Creating http server")
	httpServer, err := controller.NewHttpServer(app)
	if err != nil {
		panic(err)
	}

	logger.Infoln("Starting http server")
	if err := httpServer.Start(); err != nil {
		panic(err)
	}

	<-ctx.Done()
}
