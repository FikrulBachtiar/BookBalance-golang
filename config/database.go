package config

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"
)

var debe *sql.DB

type ConnectionDB struct {
	DBHost     string
    DBPort     int
    DBName     string
    DBUser     string
    DBPassword string
    DBmode string
}

func InitDB(con *ConnectionDB) {
	db, err := sql.Open("postgres", fmt.Sprintf(`host=%s port=%d user=%s password=%s dbname=%s sslmode=%s`, con.DBHost, con.DBPort, con.DBUser, con.DBPassword, con.DBName, con.DBmode))
	if err != nil {
		panic(fmt.Sprintf("DB Error Open => %s", err))
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	err = db.Ping()
	if err != nil {
		panic(fmt.Sprintf("Connection DB Error => %s", err))
	}
}

func DBConn() *sql.DB {
	return debe
}