package api

import (
	"testing"
	"practicas_go/internal/errors"
)

func TestNewClient(t *testing.T) {
	client := NewClient()
	if client == nil {
		t.Fatal("NewClient() returned nil")
	}
	if client.httpClient == nil {
		t.Fatal("httpClient is nil")
	}
}

func TestAnalyze_InvalidDomain(t *testing.T) {
	client := NewClient()
	_, err := client.Analyze("")
	
	if err == nil {
		t.Fatal("Expected error for empty domain")
	}
	
	if _, ok := err.(*errors.DomainError); !ok {
		t.Fatalf("Expected DomainError, got %T", err)
	}
}

func TestGetEndpointData_InvalidInput(t *testing.T) {
	client := NewClient()
	
	tests := []struct {
		name      string
		host      string
		ipAddress string
	}{
		{"empty host", "", "1.1.1.1"},
		{"empty ip", "example.com", ""},
		{"both empty", "", ""},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := client.GetEndpointData(tt.host, tt.ipAddress)
			if err == nil {
				t.Fatal("Expected error for invalid input")
			}
			if _, ok := err.(*errors.DomainError); !ok {
				t.Fatalf("Expected DomainError, got %T", err)
			}
		})
	}
}

