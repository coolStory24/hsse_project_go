package errors

import "fmt"

type ServiceBadRequestError struct {
	Message string
	Details string
}

func (e *ServiceBadRequestError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("%s: %s", e.Message, e.Details)
	}
	return e.Message
}

func NewServiceBadRequestError(message string, details string) *ServiceBadRequestError {
	return &ServiceBadRequestError{
		Message: message,
		Details: details,
	}
}
