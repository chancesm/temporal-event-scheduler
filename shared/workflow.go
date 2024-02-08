package shared

import (
	"time"

	"go.temporal.io/sdk/workflow"
)

func GenerateScheduledEvent(ctx workflow.Context, params RequestScheduledEventParams) error {
	logger := workflow.GetLogger(ctx)
	options := workflow.ActivityOptions{
		StartToCloseTimeout: time.Second * 30,
	}
	ctx = workflow.WithActivityOptions(ctx, options)
	var activityResult string
	err := workflow.ExecuteActivity(ctx, EmitScheduledEvent, params).Get(ctx, &activityResult)
	if err != nil {
		logger.Error("Activity Failed", "error", err)
	}
	return nil
}
