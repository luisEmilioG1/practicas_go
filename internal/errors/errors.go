package errors

import "fmt"

var (
	ErrInvalidDomain    = &DomainError{Message: "invalid domain provided"}
	ErrAPIConnection    = &APIError{Message: "failed to connect to SSL Labs API"}
	ErrAPIResponse      = &APIError{Message: "invalid response from SSL Labs API"}
	ErrAnalysisNotReady = &AnalysisError{Message: "analysis not ready yet, please wait"}
	ErrNoEndpoints      = &AnalysisError{Message: "no endpoints found for this domain"}
)

type DomainError struct {
	Message string
}

func (e *DomainError) Error() string {
	return e.Message
}

type APIError struct {
	Message string
	Err     error
}

func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *APIError) Unwrap() error {
	return e.Err
}

type AnalysisError struct {
	Message string
}

func (e *AnalysisError) Error() string {
	return e.Message
}
