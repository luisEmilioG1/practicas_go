package handlers

import (
	"html/template"
	"net/http"

	"practicas_go/internal/analyzer"
	"practicas_go/internal/errors"
)

type Handler struct {
	analyzer *analyzer.Analyzer
	templates *template.Template
}

func NewHandler(analyzer *analyzer.Analyzer) (*Handler, error) {
	templates, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}

	return &Handler{
		analyzer:  analyzer,
		templates: templates,
	}, nil
}

func (h *Handler) ShowForm(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "index.html", nil); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) AnalyzeDomain(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	host := r.FormValue("domain")
	if host == "" {
		h.renderError(w, "Domain is required")
		return
	}

	result, err := h.analyzer.AnalyzeDomain(host)
	if err != nil {
		h.handleError(w, err)
		return
	}

	if err := h.templates.ExecuteTemplate(w, "result.html", result); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (h *Handler) handleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *errors.DomainError:
		h.renderError(w, e.Message)
	case *errors.APIError:
		h.renderError(w, "Failed to connect to SSL Labs API. Please try again later.")
	case *errors.AnalysisError:
		h.renderError(w, e.Message)
	default:
		h.renderError(w, "An unexpected error occurred")
	}
}

func (h *Handler) renderError(w http.ResponseWriter, message string) {
	data := map[string]interface{}{
		"Error": message,
	}
	if err := h.templates.ExecuteTemplate(w, "error.html", data); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
