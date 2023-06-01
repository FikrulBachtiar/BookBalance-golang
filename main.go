package main

import (
	config "bookbalance/app/configs"
	routes "bookbalance/routes"
	"context"
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

func main() {
	
	// Env data
	envFile := ".env"
	err := godotenv.Load(envFile)
	if err != nil {
		fmt.Printf("Error loading :%s file", envFile)
	}

	appPort := os.Getenv("APP_PORT")

	DBHost := os.Getenv("DB_HOST")
	DBName := os.Getenv("DB_NAME")
	DBUser := os.Getenv("DB_USER")
	DBPassword := os.Getenv("DB_PASSWORD")
	DBmode := os.Getenv("DB_SSL_MODE")
	DBPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		panic(fmt.Sprintf("Failed to convert string to int [DB] => %s", err))
	}

	RedisHost := os.Getenv("REDIS_HOST");
	RedisPassword := os.Getenv("REDIS_PASSWORD");
	RedisDB, err := strconv.Atoi(os.Getenv("REDIS_DB"));
	if err != nil {
		panic(fmt.Sprintf("Failed to convert string to int [REDIS] => %s", err))
	}

	configDB := &config.ConnectionDB{
		DBHost: DBHost,
		DBPort: DBPort,
		DBName: DBName,
		DBUser: DBUser,
		DBPassword: DBPassword,
		DBmode: DBmode,
	}

	configRedis := &config.ConnectionRedis{
		Addr: RedisHost,
		Password: RedisPassword,
		DB: RedisDB,
	}

	// Routes
	db := config.InitDB(configDB);
	redis := config.InitRedis(context.Background(), configRedis);
	defer db.Close()
	app := routes.InitRoutes(db, redis);

	app.Start(fmt.Sprintf(":%s", appPort));
}