package main

import (
	"log"
	"myapp/models"
	"myapp/router"
	"os"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func main() {
	// set logrus
	logrus.SetLevel(logrus.TraceLevel)
	// set database
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	dsn := os.Getenv(("DATABASE_URL"))
	if dsn == "" {
		log.Fatal("DATABASE_URL is not set")
	}

	_, err := models.InitDB(dsn)
	if err != nil {
		logrus.Fatalf("Failed to connect to database: %v", err)
	}

	r := router.SetupRouter()
	r.Run()
}
