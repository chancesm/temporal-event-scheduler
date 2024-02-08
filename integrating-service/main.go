package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/chancesm/temporal-event-scheduler/shared"
	"github.com/google/uuid"
	kafka "github.com/segmentio/kafka-go"
)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(kafkaURL),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
}

func main() {
	// get kafka writer using environment variables.
	kafkaURL := "localhost:9092"
	topic := "scheduler"

	writer := newKafkaWriter(kafkaURL, topic)
	defer writer.Close()
	fmt.Println("start producing ... !!")
	key := "requestScheduledEvent"

	scheduledEvent := shared.RequestScheduledEventParams{
		RequestId: fmt.Sprint(uuid.New()),
		Schedule:  "every 30s",
		Event: shared.ScheduledEvent{
			Data:      "Some Value",
			RequestId: fmt.Sprint(uuid.New()),
		},
	}

	jsonbytes, _ := json.Marshal(scheduledEvent)
	msg := kafka.Message{
		Key:   []byte(key),
		Value: jsonbytes,
	}
	err := writer.WriteMessages(context.Background(), msg)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("produced", key)
	}

	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "scheduler",
		MaxBytes: 10e6, // 10MB
	})
	defer r.Close()

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			break
		}
		if string(m.Key) == "scheduledEvent" {
			fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
		} else {
			fmt.Println("ignored message")
		}

	}
}
