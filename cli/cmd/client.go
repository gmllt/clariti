package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

// API client structure
type APIClient struct {
	BaseURL    string
	Username   string
	Password   string
	HTTPClient *http.Client
}

// NewAPIClient creates a new API client
func NewAPIClient(baseURL, username, password string) *APIClient {
	return &APIClient{
		BaseURL:    baseURL,
		Username:   username,
		Password:   password,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// makeRequest makes an HTTP request with basic auth
func (c *APIClient) makeRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	var reqBody io.Reader
	var requestBodyString string

	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
		requestBodyString = string(jsonBody)
	}

	url := c.BaseURL + endpoint

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	if c.Username != "" && c.Password != "" {
		req.SetBasicAuth(c.Username, c.Password)
	}

	// Detailed REQUEST trace
	traceRequest(method, url, req.Header, requestBodyString)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Read response body for tracing (we need to create a new reader for the caller)
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	resp.Body.Close()

	// Detailed RESPONSE trace
	traceResponse(resp.StatusCode, resp.Header, string(respBody))

	// Create new response body reader for the caller
	resp.Body = io.NopCloser(bytes.NewBuffer(respBody))

	return resp, nil
}

// trace outputs trace information if tracing is enabled
func trace(format string, args ...interface{}) {
	if os.Getenv("CLARITI_TRACE_ENABLED") != "true" {
		return
	}

	message := fmt.Sprintf("[TRACE] "+format, args...)

	traceFile := os.Getenv("CLARITI_TRACE_FILE")
	if traceFile != "" {
		file, err := os.OpenFile(traceFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open trace file: %v\n", err)
			return
		}
		defer file.Close()
		fmt.Fprintf(file, "%s %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
	} else {
		fmt.Fprintf(os.Stdout, "%s %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
	}
}

// traceRequest outputs detailed request information
func traceRequest(method, url string, headers http.Header, body string) {
	if os.Getenv("CLARITI_TRACE_ENABLED") != "true" {
		return
	}

	message := fmt.Sprintf("[REQUEST] %s %s", method, url)
	writeTrace(message)

	// Trace headers (excluding sensitive auth)
	for key, values := range headers {
		if key != "Authorization" {
			for _, value := range values {
				writeTrace(fmt.Sprintf("[REQUEST] %s: %s", key, value))
			}
		} else {
			writeTrace(fmt.Sprintf("[REQUEST] %s: [HIDDEN]", key))
		}
	}

	// Trace body if present
	if body != "" {
		writeTrace(fmt.Sprintf("[REQUEST] Body: %s", body))
	}
}

// traceResponse outputs detailed response information
func traceResponse(statusCode int, headers http.Header, body string) {
	if os.Getenv("CLARITI_TRACE_ENABLED") != "true" {
		return
	}

	message := fmt.Sprintf("[RESPONSE] Status: %d", statusCode)
	writeTrace(message)

	// Trace important headers
	for key, values := range headers {
		if key == "Content-Type" || key == "Content-Length" {
			for _, value := range values {
				writeTrace(fmt.Sprintf("[RESPONSE] %s: %s", key, value))
			}
		}
	}

	// Trace body
	if body != "" {
		writeTrace(fmt.Sprintf("[RESPONSE] Body: %s", body))
	}
}

// writeTrace writes trace message to output
func writeTrace(message string) {
	traceFile := os.Getenv("CLARITI_TRACE_FILE")
	if traceFile != "" {
		file, err := os.OpenFile(traceFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Failed to open trace file: %v\n", err)
			return
		}
		defer file.Close()
		fmt.Fprintf(file, "%s %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
	} else {
		fmt.Fprintf(os.Stdout, "%s %s\n", time.Now().Format("2006-01-02 15:04:05"), message)
	}
}

// getAPIClient returns a configured API client
func getAPIClient() *APIClient {
	// Try to get target configuration first
	target, err := getCurrentTarget()
	if err == nil {
		// Use target configuration
		return NewAPIClient(target.URL, target.Username, target.Password)
	}

	// Fallback to flags/environment variables
	return NewAPIClient(serverURL, username, password)
}
