package shared

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func GenerateScheduledEvent(ctx workflow.Context, params RequestScheduledEventParams) error {
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 30,
	}
	ctx = workflow.WithActivityOptions(ctx, options)

	workflow.ExecuteActivity(ctx, EmitScheduledEvent, params)

	return nil
}
