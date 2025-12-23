package analyzer

import (
	"practicas_go/internal/api"
	"practicas_go/internal/errors"
	"practicas_go/internal/models"
	"strings"
)

type Analyzer struct {
	client *api.Client
}

func NewAnalyzer(client *api.Client) *Analyzer {
	return &Analyzer{
		client: client,
	}
}

func (a *Analyzer) AnalyzeDomain(host string) (*models.AnalysisResult, error) {
	if host == "" {
		return nil, errors.ErrInvalidDomain
	}

	analyzeResp, err := a.client.WaitForAnalysis(host)
	if err != nil {
		return nil, err
	}

	if len(analyzeResp.Endpoints) == 0 {
		return nil, errors.ErrNoEndpoints
	}

	var endpointResults []models.EndpointResult
	var overallGrade string
	var hasWarnings bool
	var isExceptional bool
	grades := make(map[string]int)

	for _, endpoint := range analyzeResp.Endpoints {
		endpointData, err := a.client.GetEndpointData(host, endpoint.IPAddress)
		if err != nil {
			continue
		}

		var protocols []string
		for _, protocol := range endpointData.Protocols {
			if protocol.Version != "" {
				protocols = append(protocols, protocol.Version)
			}
		}

		endpointResult := models.EndpointResult{
			IPAddress:    endpointData.IPAddress,
			Grade:        endpointData.Grade,
			Protocols:    protocols,
			HasWarnings:  endpointData.HasWarnings,
			IsExceptional: endpointData.IsExceptional,
		}

		endpointResults = append(endpointResults, endpointResult)

		if endpointData.Grade != "" {
			grades[endpointData.Grade]++
		}
		if endpointData.HasWarnings {
			hasWarnings = true
		}
		if endpointData.IsExceptional {
			isExceptional = true
		}
	}

	overallGrade = calculateOverallGrade(grades)
	recommendations := generateRecommendations(endpointResults, hasWarnings, overallGrade)

	return &models.AnalysisResult{
		Domain:          host,
		OverallGrade:    overallGrade,
		HasWarnings:     hasWarnings,
		IsExceptional:   isExceptional,
		Endpoints:       endpointResults,
		Recommendations: recommendations,
	}, nil
}

func calculateOverallGrade(grades map[string]int) string {
	if len(grades) == 0 {
		return "N/A"
	}

	gradeOrder := map[string]int{
		"A+": 7, "A": 6, "A-": 5,
		"B": 4, "B-": 3,
		"C": 2, "C-": 1,
		"D": 0, "E": -1, "F": -2,
		"T": -3, "M": -4,
	}

	lowestScore := 10
	lowestGrade := "A+"

	for grade, count := range grades {
		if count > 0 {
			score, exists := gradeOrder[grade]
			if !exists {
				continue
			}
			if score < lowestScore {
				lowestScore = score
				lowestGrade = grade
			}
		}
	}

	return lowestGrade
}

func generateRecommendations(endpoints []models.EndpointResult, hasWarnings bool, overallGrade string) []string {
	var recommendations []string

	if overallGrade == "F" || overallGrade == "T" || overallGrade == "M" {
		recommendations = append(recommendations, "CRITICAL: Immediate action required. TLS configuration is insecure.")
	} else if overallGrade == "D" || overallGrade == "E" {
		recommendations = append(recommendations, "WARNING: TLS configuration needs improvement. Security vulnerabilities detected.")
	} else if overallGrade == "C" || overallGrade == "C-" {
		recommendations = append(recommendations, "Consider upgrading TLS configuration for better security.")
	}

	if hasWarnings {
		recommendations = append(recommendations, "Some endpoints have security warnings. Review configuration details.")
	}

	protocolIssues := false
	for _, endpoint := range endpoints {
		hasTLS13 := false
		hasTLS12 := false
		for _, protocol := range endpoint.Protocols {
			if strings.Contains(protocol, "1.3") {
				hasTLS13 = true
			}
			if strings.Contains(protocol, "1.2") {
				hasTLS12 = true
			}
		}
		if !hasTLS13 && !hasTLS12 {
			protocolIssues = true
		}
	}

	if protocolIssues {
		recommendations = append(recommendations, "Some endpoints do not support modern TLS protocols (1.2 or 1.3).")
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "TLS configuration appears secure. No immediate action required.")
	}

	return recommendations
}
