package shared

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/activity"
)

func newKafkaWriter(kafkaURL, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(kafkaURL),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		AllowAutoTopicCreation: true,
	}
}

func EmitScheduledEvent(ctx context.Context, params RequestScheduledEventParams) (string, error) {
	logger := activity.GetLogger(ctx)
	logger.Info("Emit Scheduled Logger Called with params", "params", params)
	kafkaURL := "localhost:9092"
	topic := "scheduler"
	writer := newKafkaWriter(kafkaURL, topic)
	defer writer.Close()

	jsonString, _ := json.Marshal(params.Event)

	key := "scheduledEvent"
	msg := kafka.Message{
		Key:   []byte(key),
		Value: jsonString,
	}
	err := writer.WriteMessages(context.Background(), msg)
	if err != nil {
		logger.Error("Failed to Write", "error", err)
		return "", err
	} else {
		logger.Info("Successfully wrote message", "message", msg)
	}
	return "done", nil
}
