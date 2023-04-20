package main

import (
	"log"

	"github.com/ResulShamuhammedov/Kafka-Clickhouse/server/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal(err)
	}

	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "test-topic",
	})
	defer writer.Close()

	handler := handler.NewHandler(writer)

	app := fiber.New()

	app.Post("/metric", handler.HandleGetMetrics)

	err := app.Listen(":8080")
	if err != nil {
		log.Fatalf("could not listen: %s", err.Error())
	}
}