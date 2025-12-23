package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"practicas_go/internal/errors"
	"practicas_go/internal/models"
)

const (
	baseURL     = "https://api.ssllabs.com/api/v3"
	maxWaitTime = 60 * time.Second
	pollInterval = 5 * time.Second
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Analyze(host string) (*models.AnalyzeResponse, error) {
	if host == "" {
		return nil, errors.ErrInvalidDomain
	}

	url := fmt.Sprintf("%s/analyze?host=%s&all=done", baseURL, host)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, &errors.APIError{
			Message: errors.ErrAPIConnection.Message,
			Err:     err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &errors.APIError{
			Message: errors.ErrAPIResponse.Message,
			Err:     fmt.Errorf("status code: %d", resp.StatusCode),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.APIError{
			Message: errors.ErrAPIResponse.Message,
			Err:     err,
		}
	}

	var analyzeResp models.AnalyzeResponse
	if err := json.Unmarshal(body, &analyzeResp); err != nil {
		return nil, &errors.APIError{
			Message: errors.ErrAPIResponse.Message,
			Err:     err,
		}
	}

	if analyzeResp.Status == "ERROR" {
		return nil, &errors.APIError{
			Message: errors.ErrAPIResponse.Message,
			Err:     fmt.Errorf("API returned error status"),
		}
	}

	return &analyzeResp, nil
}

func (c *Client) WaitForAnalysis(host string) (*models.AnalyzeResponse, error) {
	startTime := time.Now()

	for {
		if time.Since(startTime) > maxWaitTime {
			return nil, &errors.APIError{
				Message: errors.ErrAPIResponse.Message,
				Err:     fmt.Errorf("analysis timeout"),
			}
		}

		analyzeResp, err := c.Analyze(host)
		if err != nil {
			return nil, err
		}

		if analyzeResp.Status == "READY" || analyzeResp.Status == "ERROR" {
			return analyzeResp, nil
		}

		time.Sleep(pollInterval)
	}
}

func (c *Client) GetEndpointData(host, ipAddress string) (*models.EndpointDataResponse, error) {
	if host == "" || ipAddress == "" {
		return nil, errors.ErrInvalidDomain
	}

	url := fmt.Sprintf("%s/getEndpointData?host=%s&s=%s", baseURL, host, ipAddress)
	
	resp, err := c.httpClient.Get(url)
	if err != nil {
		return nil, &errors.APIError{
			Message: errors.ErrAPIConnection.Message,
			Err:     err,
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, &errors.APIError{
			Message: errors.ErrAPIResponse.Message,
			Err:     fmt.Errorf("status code: %d", resp.StatusCode),
		}
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, &errors.APIError{
			Message: errors.ErrAPIResponse.Message,
			Err:     err,
		}
	}

	var endpointData models.EndpointDataResponse
	if err := json.Unmarshal(body, &endpointData); err != nil {
		return nil, &errors.APIError{
			Message: errors.ErrAPIResponse.Message,
			Err:     err,
		}
	}

	return &endpointData, nil
}
