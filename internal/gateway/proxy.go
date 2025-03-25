package gateway

import (
	"io"
	"net/http"

	"github.com/leo-andrei/api-gateway/config"
)

// CreateProxyHandler creates a handler function for a given route
func CreateProxyHandler(route config.Route) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Create a new request to the target URL
		client := &http.Client{}
		req, err := http.NewRequest(route.Method, route.TargetURL, r.Body)
		if err != nil {
			http.Error(w, "Error creating request", http.StatusInternalServerError)
			return
		}

		// Copy headers from the original request
		for name, values := range r.Header {
			for _, value := range values {
				req.Header.Add(name, value)
			}
		}

		// Add X-Forwarded headers
		req.Header.Add("X-Forwarded-For", r.RemoteAddr)
		req.Header.Add("X-Forwarded-Host", r.Host)
		req.Header.Add("X-Forwarded-Proto", "http") // or "https" if using TLS

		// Make the request to the target URL
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, "Error forwarding request", http.StatusBadGateway)
			return
		}
		defer resp.Body.Close()

		// Copy response headers to the client response
		for name, values := range resp.Header {
			for _, value := range values {
				w.Header().Add(name, value)
			}
		}

		// Set the status code
		w.WriteHeader(resp.StatusCode)

		// Copy the response body to the client
		io.Copy(w, resp.Body)
	}
}
