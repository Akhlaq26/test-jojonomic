package config

import (
	"context"
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
}

func NewConfig() *Config {
	err := godotenv.Load("./.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	rt := mux.NewRouter()
	kafka, err := kafka.DialLeader(context.Background(), "tcp", os.Getenv("KAFKA_URL"), os.Getenv("KAFKA_TOPIC"), 0)
	if err != nil {
		log.Fatalf("kafka connection err : %v", err.Error())
	}
	app := &Config{
		Rt:    rt,
		Kafka: kafka,
	}

	return app
}

func URL() string {
	return os.Getenv("URL")
}
