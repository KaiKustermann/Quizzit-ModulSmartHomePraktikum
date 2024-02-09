package api

import (
	"fmt"
	"net/http"
)

// Handle incoming health requests
func HealthCheckHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintf(w, "System is running...")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
