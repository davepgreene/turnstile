package http

import "net/http"

type WrappedResponseWriter struct {
	status int
	wroteHeader bool
	http.ResponseWriter
}

func NewWrappedResponseWriter(rw http.ResponseWriter) *WrappedResponseWriter {
	return &WrappedResponseWriter{ResponseWriter: rw}
}

// Give a way to get the status
func (w *WrappedResponseWriter) Status() int {
	return w.status
}

// Satisfy the http.ResponseWriter interface
func (w *WrappedResponseWriter) Header() http.Header {
	return w.ResponseWriter.Header()
}

func (w *WrappedResponseWriter) Write(data []byte) (n int, err error) {
	return w.ResponseWriter.Write(data)
}

func (w *WrappedResponseWriter) WriteHeader(statusCode int) {
	// Store the status code
	w.status = statusCode

	// Write the status code onward.
	w.ResponseWriter.WriteHeader(statusCode)
}