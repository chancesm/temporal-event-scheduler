package main

// This is a placeholder program, unrelated to the training course. You
// are not expected to run or modify this file. Having a Go program in
// the project root avoids warning messages shown when running "go get"
// from this directory. See: https://github.com/golang/go/issues/37700

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/chancesm/temporal-event-scheduler/shared"
	"github.com/gofiber/fiber/v2"
	"github.com/segmentio/kafka-go"
	"go.temporal.io/sdk/client"
)

func main() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})
	app.Post("/schedule", func(c *fiber.Ctx) error {
		var params shared.RequestScheduledEventParams
		if err := c.BodyParser(params); err != nil {
			return err
		}
		scheduleEvent(params)

		return c.JSON(struct {
			Message string
			Data    any
		}{
			Message: "Successfully Scheduled Event",
			Data:    params,
		})
	})

	go func() {
		log.Fatal(app.Listen(":3000"))
	}()

	// make a new reader that consumes from topic-A, partition 0, at offset 42
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  []string{"localhost:9092"},
		Topic:    "scheduler",
		MaxBytes: 10e6, // 10MB
	})
	r.SetOffsetAt(context.Background(), time.Now())

	go func() {
		for {
			m, err := r.ReadMessage(context.Background())
			if err != nil {
				break
			}
			if string(m.Key) == "requestScheduledEvent" {
				var params shared.RequestScheduledEventParams
				json.Unmarshal(m.Value, &params)

				scheduleEvent(params)
				fmt.Printf("message at offset %d: %s = %s\n", m.Offset, string(m.Key), string(m.Value))
			}
		}
	}()

	<-done
	log.Print("Server Stopped")
}

func scheduleEvent(params shared.RequestScheduledEventParams) {
	log.Println("This function was called")
	temporalClient, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer temporalClient.Close()

	scheduleID := params.RequestId
	workflowID := fmt.Sprintf("%s_schedule_workflow", scheduleID)
	// Create the schedule.
	_, err = temporalClient.ScheduleClient().Create(context.TODO(), client.ScheduleOptions{
		ID: scheduleID,
		Spec: client.ScheduleSpec{
			Intervals: []client.ScheduleIntervalSpec{
				{Every: time.Second * 30},
			},
		},
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  shared.GenerateScheduledEvent,
			TaskQueue: shared.SchedulerTemporalQueue,
			Args:      []interface{}{params},
		},
	})
	if err != nil {
		log.Fatalln("Unable to create Schedule", err)
	}
}
