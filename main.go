package main

import (
	"log"
	"myapp/models"
	"myapp/router"
	"os"
	"time"

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

	// retry init database connection
	var err error
	for i := 0; i < 10; i++ {

		_, err = models.InitDB(dsn)
		if err == nil {
			logrus.Info("Successfully connected to the database.")
			break
		}
		logrus.Errorf("Failed to connect to database: %v. Retrying in 2 seconds...", err)
		time.Sleep(2 * time.Second)
	}

	if err != nil {
		logrus.Fatalf("Failed to connect to database after retries: %v", err)
	}

	// _, err := models.InitDB(dsn)
	// if err != nil {
	// 	logrus.Fatalf("Failed to connect to database: %v", err)
	// }

	r := router.SetupRouter()
	r.Run()
}
