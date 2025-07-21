package models

type BulkOperationResult struct {
	SuccessCount int     `json:"successCount"`
	FailureCount int     `json:"failureCount"`
	Errors       []error `json:"errors,omitempty"`
}
