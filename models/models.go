package models

import "time"

type KeyValue struct {
	Value      string
	ExpiryTime *time.Time
}

type Queue struct {
	Items []string
}
type CommandRequest struct {
	Command string `json:"command"`
}

type CommandResponse struct {
	Value string `json:"value,omitempty"`
	Error string `json:"error,omitempty"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
