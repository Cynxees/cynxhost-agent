package middleware

import (
	"bytes"
	"cynxhostagent/internal/helper"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware is a middleware that logs details of the HTTP request and response.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Log the incoming request details
		startTime := time.Now()
		log.Printf("Incoming Request: %s %s\n", r.Method, r.URL.Path)
		log.Printf("Headers: %v\n", r.Header)
		log.Printf("Remote Addr: %s\n", r.RemoteAddr)
		log.Printf("Client IP: %s\n", helper.GetClientIP(r))

		// Read the body of the request (if possible)
		var requestBody string
		if r.Body != nil {
			bodyBytes, err := ioutil.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading request body: %v\n", err)
			} else {
				requestBody = string(bodyBytes)
			}
			// Re-assign the body to allow the next handler to read it
			r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		if requestBody != "" {
			log.Printf("Request Body: %s\n", requestBody)
		}

		// Create a ResponseWriter wrapper to capture the status code and response body
		lrw := &LoggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK, body: &bytes.Buffer{}}
		next.ServeHTTP(lrw, r)

		// Log the response details
		log.Printf("Response Status: %d\n", lrw.statusCode)
		log.Printf("Response Body: %s\n", lrw.body.String())
		log.Printf("Request processed in %s\n", time.Since(startTime))
	})
}

// LoggingResponseWriter wraps the standard http.ResponseWriter to capture status code and response body
type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (lrw *LoggingResponseWriter) WriteHeader(statusCode int) {
	lrw.statusCode = statusCode
	lrw.ResponseWriter.WriteHeader(statusCode)
}

func (lrw *LoggingResponseWriter) Write(p []byte) (n int, err error) {
	lrw.body.Write(p)
	return lrw.ResponseWriter.Write(p)
}
