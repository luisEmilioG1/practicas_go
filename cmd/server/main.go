package main

import (
	"fmt"
	"log"
	"net/http"

	"practicas_go/internal/analyzer"
	"practicas_go/internal/api"
	"practicas_go/internal/handlers"
)

func main() {
	client := api.NewClient()
	analyzer := analyzer.NewAnalyzer(client)
	
	handler, err := handlers.NewHandler(analyzer)
	if err != nil {
		log.Fatalf("Failed to initialize handler: %v", err)
	}

	http.HandleFunc("/", handler.ShowForm)
	http.HandleFunc("/analyze", handler.AnalyzeDomain)

	fmt.Println("Server started at http://localhost:8888")
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
