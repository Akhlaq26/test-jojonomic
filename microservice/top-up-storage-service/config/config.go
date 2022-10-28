package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	kafka "github.com/segmentio/kafka-go"
)

type Config struct {
	Rt *mux.Router
	Db *sql.DB
	Kr *kafka.Reader
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
	brokers := strings.Split(os.Getenv("KAFKA_URL"), ",")
	kafkaReader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  brokers,
		GroupID:  os.Getenv("KAFKA_GROUP_ID"),
		Topic:    os.Getenv("KAFKA_TOPIC"),
		MinBytes: 0,    // 1KB
		MaxBytes: 10e6, // 10MB
	})
	app := &Config{
		Db: db,
		Rt: rt,
		Kr: kafkaReader,
	}

	return app
}

func URL() string {
	return os.Getenv("URL")
}
