package main

import (
	"bookbalance/config"
	"bookbalance/routes"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	envFile := ".env"
	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Printf("Error loading :%s file", envFile)
	}

	// Database
	DBHost := os.Getenv("DB_HOST")
	DBName := os.Getenv("DB_NAME")
	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBmode := os.Getenv("DB_SSL_MODE")
	DBPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic(fmt.Sprintf("Failed to convert string to int => %s", err))
	}

	configDB := &config.ConnectionDB{
		DBHost: DBHost,
		DBPort: DBPort,
		DBName: DBName,
		DBUser: DBUser,
		DBPassword: DBPassword,
		DBmode: DBmode,
	}

	config.InitDB(configDB)

	// Routes
	appPort := os.Getenv("APP_PORT")
	app := routes.InitRoutes()

	app.Start(fmt.Sprintf(":%s", appPort));
}