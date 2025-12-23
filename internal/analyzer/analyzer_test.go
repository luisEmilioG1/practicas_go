package analyzer

import (
	"testing"
	"practicas_go/internal/api"
	"practicas_go/internal/errors"
	"practicas_go/internal/models"
)

func TestNewAnalyzer(t *testing.T) {
	client := api.NewClient()
	analyzer := NewAnalyzer(client)
	
	if analyzer == nil {
		t.Fatal("NewAnalyzer() returned nil")
	}
	if analyzer.client == nil {
		t.Fatal("client is nil")
	}
}

func TestAnalyzeDomain_InvalidDomain(t *testing.T) {
	client := api.NewClient()
	analyzer := NewAnalyzer(client)
	
	_, err := analyzer.AnalyzeDomain("")
	
	if err == nil {
		t.Fatal("Expected error for empty domain")
	}
	
	if _, ok := err.(*errors.DomainError); !ok {
		t.Fatalf("Expected DomainError, got %T", err)
	}
}

func TestCalculateOverallGrade(t *testing.T) {
	tests := []struct {
		name   string
		grades map[string]int
		want   string
	}{
		{
			name:   "single A grade",
			grades: map[string]int{"A": 1},
			want:   "A",
		},
		{
			name:   "multiple grades, lowest is F",
			grades: map[string]int{"A": 2, "F": 1},
			want:   "F",
		},
		{
			name:   "empty grades",
			grades: map[string]int{},
			want:   "N/A",
		},
		{
			name:   "B and C grades",
			grades: map[string]int{"B": 1, "C": 1},
			want:   "C",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := calculateOverallGrade(tt.grades)
			if got != tt.want {
				t.Errorf("calculateOverallGrade() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGenerateRecommendations(t *testing.T) {
	tests := []struct {
		name         string
		endpoints    []models.EndpointResult
		hasWarnings  bool
		overallGrade string
		wantMin      int
	}{
		{
			name:         "F grade should have recommendations",
			endpoints:    []models.EndpointResult{},
			hasWarnings:  false,
			overallGrade: "F",
			wantMin:      1,
		},
		{
			name:         "has warnings should have recommendations",
			endpoints:    []models.EndpointResult{},
			hasWarnings:  true,
			overallGrade: "A",
			wantMin:      1,
		},
		{
			name:         "A grade without warnings",
			endpoints:    []models.EndpointResult{},
			hasWarnings:  false,
			overallGrade: "A",
			wantMin:      1,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := generateRecommendations(tt.endpoints, tt.hasWarnings, tt.overallGrade)
			if len(got) < tt.wantMin {
				t.Errorf("generateRecommendations() returned %d recommendations, want at least %d", len(got), tt.wantMin)
			}
		})
	}
}

