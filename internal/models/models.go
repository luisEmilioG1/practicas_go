package models

type AnalysisRequest struct {
	Host string `json:"host"`
}

type AnalyzeResponse struct {
	Status    string      `json:"status"`
	Host      string      `json:"host"`
	Endpoints []Endpoint  `json:"endpoints"`
}

type Endpoint struct {
	IPAddress string `json:"ipAddress"`
	Grade     string `json:"grade"`
}

type EndpointDataResponse struct {
	IPAddress         string   `json:"ipAddress"`
	ServerName        string   `json:"serverName"`
	Grade             string   `json:"grade"`
	GradeTrustIgnored string   `json:"gradeTrustIgnored"`
	HasWarnings       bool     `json:"hasWarnings"`
	IsExceptional     bool     `json:"isExceptional"`
	Protocols         []Protocol `json:"protocols"`
	Suites            []Suite    `json:"suites"`
}

type Protocol struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Version  string `json:"version"`
	V2Suites []string `json:"v2Suites"`
}

type Suite struct {
	List   int    `json:"list"`
	Strength int  `json:"strength"`
	Name   string `json:"name"`
	CipherStrength int `json:"cipherStrength"`
}

type AnalysisResult struct {
	Domain           string
	OverallGrade     string
	HasWarnings      bool
	IsExceptional    bool
	Endpoints        []EndpointResult
	Recommendations  []string
}

type EndpointResult struct {
	IPAddress    string
	Grade        string
	Protocols    []string
	HasWarnings  bool
	IsExceptional bool
}
