package responsewriter

import "net/http"

// ResponseWriter is a wrapper for http.ResponseWriter that captures the status code and response size
type ResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

// NewResponseWriter creates a new ResponseWriter
func NewResponseWriter(w http.ResponseWriter) *ResponseWriter {
	return &ResponseWriter{
		ResponseWriter: w,
		statusCode:     http.StatusOK, // Default status code
	}
}

// WriteHeader captures the status code
func (rw *ResponseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// Write captures the response size
func (rw *ResponseWriter) Write(b []byte) (int, error) {
	size, err := rw.ResponseWriter.Write(b)
	rw.size += size
	return size, err
}

// StatusCode returns the HTTP status code
func (rw *ResponseWriter) StatusCode() int {
	return rw.statusCode
}

// Size returns the total size of the response
func (rw *ResponseWriter) Size() int {
	return rw.size
}
