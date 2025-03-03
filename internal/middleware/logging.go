package middleware

import (
	"bytes"
	"io"
	"time"

	"ecom-go/pkg/logger"
	"github.com/gin-gonic/gin"
)

// responseWriter is a wrapper for gin.ResponseWriter to capture the response body
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write captures the response body
func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// Logger is a middleware that logs HTTP requests and responses
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()

		// Read the request body
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create a response body buffer
		responseBodyBuffer := &bytes.Buffer{}

		// Replace the writer with our custom one
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           responseBodyBuffer,
		}
		c.Writer = writer

		// Process request
		c.Next()

		// Stop timer
		latency := time.Since(start)

		// Determine status for logging
		status := c.Writer.Status()

		// Determine log level based on status code
		var logFn func(string, ...interface{})
		if status >= 500 {
			logFn = logger.Error
		} else if status >= 400 {
			logFn = logger.Warn
		} else {
			logFn = logger.Info
		}

		// Truncate request and response bodies if they're too large
		const maxBodyLogSize = 1024 // 1KB

		requestBodyLog := requestBody
		if len(requestBodyLog) > maxBodyLogSize {
			requestBodyLog = requestBodyLog[:maxBodyLogSize]
		}

		responseBodyLog := responseBodyBuffer.Bytes()
		if len(responseBodyLog) > maxBodyLogSize {
			responseBodyLog = responseBodyLog[:maxBodyLogSize]
		}

		// Log request details
		logFn("HTTP Request",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"status", status,
			"latency", latency,
			"client_ip", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
			"request_body", string(requestBodyLog),
			"response_body", string(responseBodyLog),
		)
	}
}
