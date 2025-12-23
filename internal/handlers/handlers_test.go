package handlers

import (
	"html/template"
	"net/http"
	"net/http/httptest"
	"testing"
	"practicas_go/internal/analyzer"
	"practicas_go/internal/api"
)

func createTestHandler(analyzer *analyzer.Analyzer) *Handler {
	templates := template.Must(template.New("index.html").Parse(`<html><body>Form</body></html>`))
	template.Must(templates.New("result.html").Parse(`<html><body>Result</body></html>`))
	template.Must(templates.New("error.html").Parse(`<html><body>Error: {{.Error}}</body></html>`))
	
	return &Handler{
		analyzer:  analyzer,
		templates: templates,
	}
}

func TestHandler_ShowForm(t *testing.T) {
	client := api.NewClient()
	analyzer := analyzer.NewAnalyzer(client)
	handler := createTestHandler(analyzer)

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	handler.ShowForm(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200, got %d", w.Code)
	}
}

func TestHandler_ShowForm_MethodNotAllowed(t *testing.T) {
	client := api.NewClient()
	analyzer := analyzer.NewAnalyzer(client)
	handler := createTestHandler(analyzer)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	handler.ShowForm(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandler_AnalyzeDomain_MethodNotAllowed(t *testing.T) {
	client := api.NewClient()
	analyzer := analyzer.NewAnalyzer(client)
	handler := createTestHandler(analyzer)

	req := httptest.NewRequest(http.MethodGet, "/analyze", nil)
	w := httptest.NewRecorder()

	handler.AnalyzeDomain(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected status 405, got %d", w.Code)
	}
}

func TestHandler_AnalyzeDomain_EmptyDomain(t *testing.T) {
	client := api.NewClient()
	analyzer := analyzer.NewAnalyzer(client)
	handler := createTestHandler(analyzer)

	req := httptest.NewRequest(http.MethodPost, "/analyze", nil)
	w := httptest.NewRecorder()

	handler.AnalyzeDomain(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected status 200 (error page), got %d", w.Code)
	}
}

