// Package healthapi defines endpoints to handle requests related to System Health
package healthapi

import (
	"fmt"
	"net/http"
)

// Handle incoming health requests on GET or return METHOD_NOT_ALLOWED
func HealthCheckHttp(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		fmt.Fprintf(w, "System is running...")
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}
