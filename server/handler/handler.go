package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Shopify/sarama"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
)

type HandlerV1 struct {
	KafkaWriter *kafka.Writer
}

func NewHandlerV1(writer *kafka.Writer) *HandlerV1 {
	return &HandlerV1{writer}
}

type HandlerV2 struct {
	SaramaProducer sarama.SyncProducer
}

func NewHandlerV2(producer sarama.SyncProducer) *HandlerV2 {
	return &HandlerV2{producer}
}

type RequestBody struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

func (h *HandlerV1) HandleGetMetrics(c *fiber.Ctx) error {
	body := new(RequestBody)
	if err := c.BodyParser(body); err != nil {
		return newErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	err = h.KafkaWriter.WriteMessages(c.UserContext(), kafka.Message{
		Value: jsonData,
	})
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}
	return c.SendStatus(http.StatusOK)
}

func (h *HandlerV1) HandleInsert(c *fiber.Ctx) error {
	for i := 1; i <= 1000000; i++ {
		jsonData, _ := json.Marshal(RequestBody{
			Name: strconv.Itoa(i),
			Age:  i,
		})
		err := h.KafkaWriter.WriteMessages(c.UserContext(), kafka.Message{
			Value: jsonData,
		})
		if err != nil {
			return newErrorResponse(c, http.StatusInternalServerError, err.Error())
		}
	}
	return c.SendStatus(http.StatusOK)
}

func (h *HandlerV2) HandleGetMetrics(c *fiber.Ctx) error {
	body := new(RequestBody)
	if err := c.BodyParser(body); err != nil {
		return newErrorResponse(c, http.StatusBadRequest, err.Error())
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	partition, offset, err := h.SaramaProducer.SendMessage(&sarama.ProducerMessage{
		Topic: "test-topic",
		Value: sarama.StringEncoder(string(jsonData)),
	})
	fmt.Printf("Message sent to partition %d at offset %d\n", partition, offset)

	if err != nil {
		return newErrorResponse(c, http.StatusInternalServerError, err.Error())
	}

	return c.SendStatus(http.StatusOK)
}
