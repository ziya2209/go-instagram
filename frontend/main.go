package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// CORS middleware to handle cross-origin requests
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Add CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Custom file server handler that adds proper MIME types
func createFileServer() http.Handler {
	fileServer := http.FileServer(http.Dir("."))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set correct MIME types for common file extensions
		ext := filepath.Ext(r.URL.Path)
		switch ext {
		case ".css":
			w.Header().Set("Content-Type", "text/css")
		case ".js":
			w.Header().Set("Content-Type", "application/javascript")
		case ".html":
			w.Header().Set("Content-Type", "text/html")
		}

		fileServer.ServeHTTP(w, r)
	})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Create file server with CORS support
	fileServer := createFileServer()
	handler := corsMiddleware(fileServer)

	// Setup the server
	mux := http.NewServeMux()
	mux.Handle("/", handler)

	fmt.Printf("Frontend server running at http://localhost:%s\n", port)
	fmt.Printf("Serving files from current directory\n")
	fmt.Println("Press Ctrl+C to stop the server")

	log.Fatal(http.ListenAndServe(":"+port, mux))
}
