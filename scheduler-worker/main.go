package main

import (
	"log"

	"github.com/chancesm/temporal-event-scheduler/shared"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Unable to create client", err)
	}
	defer c.Close()

	w := worker.New(c, shared.SchedulerTemporalQueue, worker.Options{})

	w.RegisterWorkflow(shared.GenerateScheduledEvent)
	w.RegisterActivity(shared.EmitScheduledEvent)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Unable to start worker", err)
	}
}
