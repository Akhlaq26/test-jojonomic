package config

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
	kafka "github.com/segmentio/kafka-go"
)

type Config struct {
	Rt    *mux.Router
	Kafka *kafka.Conn
	Db    *sql.DB
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
	kafka, err := kafka.DialLeader(context.Background(), "tcp", os.Getenv("KAFKA_URL"), os.Getenv("KAFKA_TOPIC"), 0)
	if err != nil {
		log.Fatal(err.Error())
	}
	app := &Config{
		Rt:    rt,
		Kafka: kafka,
		Db:    db,
	}

	return app
}

func URL() string {
	return os.Getenv("URL")
}
