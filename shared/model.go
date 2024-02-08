package shared

type RequestScheduledEventParams struct {
	Event     ScheduledEvent
	Schedule  string
	RequestId string
}

type ScheduledEvent struct {
	Data      interface{}
	RequestId string
}
