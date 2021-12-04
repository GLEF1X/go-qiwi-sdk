package core

import (
	"fmt"
	"time"
)

type HTTPError struct {
	Status int
}

func (e HTTPError) Error() string {
	return fmt.Sprintf("Http %d", e.Status)
}

type APIError struct {
	DateTime    time.Time `json:"dateTime"`
	ServiceName string    `json:"serviceName"`
	ErrorCode   string    `json:"errorCode"`
	Description string    `json:"description"`
	UserMessage string    `json:"userMessage"`
	TraceID     string    `json:"traceId"`
}

func (e APIError) Error() string {
	return fmt.Sprint(e.ErrorCode, " ", e.Description)
}
