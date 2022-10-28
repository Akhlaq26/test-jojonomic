package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type Config struct {
	Rt *mux.Router
	Db *sql.DB
}

func NewConfig() *Config {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	connectionString := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"))

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}
	rt := mux.NewRouter()
	app := &Config{
		Db: db,
		Rt: rt,
	}

	return app
}

func URL() string {
	return os.Getenv("URL")
}
