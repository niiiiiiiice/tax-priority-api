package dtos

import "time"

type CommandResult struct {
	ID        string    `json:"id,omitempty"`
	Success   bool      `json:"success"`
	Message   string    `json:"message,omitempty"`
	Error     string    `json:"error,omitempty"`
	CreatedAt time.Time `json:"createdAt,omitempty"`
	UpdatedAt time.Time `json:"updatedAt,omitempty"`
}

type BatchCommandResult struct {
	SuccessCount int             `json:"successCount"`
	FailureCount int             `json:"failureCount"`
	Results      []CommandResult `json:"results"`
	Errors       []string        `json:"errors,omitempty"`
}
