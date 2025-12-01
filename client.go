package infura

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	// BaseURL is the base URL for Infura Gas API
	BaseURL = "https://gas.api.infura.io"
	// DefaultTimeout is the default HTTP client timeout
	DefaultTimeout = 30 * time.Second
)

// Client represents the Infura Gas API client
type Client struct {
	apiKey       string
	apiKeySecret string
	baseURL      string
	httpClient   *http.Client
	debug        bool
}

// NewClient creates a new Infura Gas API client
// If apiKeySecret is empty, only API Key authentication will be used (API Key in URL path)
// If apiKeySecret is provided, Basic Auth will be used (API Key + Secret)
func NewClient(apiKey, apiKeySecret string) *Client {
	return &Client{
		apiKey:       apiKey,
		apiKeySecret: apiKeySecret,
		baseURL:      BaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// NewClientWithAPIKey creates a new client using only API Key (no secret)
// This uses the URL path authentication method: /v3/{apiKey}/networks/{chainId}/suggestedGasFees
func NewClientWithAPIKey(apiKey string) *Client {
	return &Client{
		apiKey:       apiKey,
		apiKeySecret: "",
		baseURL:      BaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}
}

// NewClientWithOptions creates a new client with custom options
// If apiKeySecret is empty, only API Key authentication will be used
func NewClientWithOptions(apiKey, apiKeySecret string, opts ...ClientOption) *Client {
	client := &Client{
		apiKey:       apiKey,
		apiKeySecret: apiKeySecret,
		baseURL:      BaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// NewClientWithAPIKeyAndOptions creates a new client with only API Key and custom options
func NewClientWithAPIKeyAndOptions(apiKey string, opts ...ClientOption) *Client {
	client := &Client{
		apiKey:       apiKey,
		apiKeySecret: "",
		baseURL:      BaseURL,
		httpClient: &http.Client{
			Timeout: DefaultTimeout,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// ClientOption is a function that configures a Client
type ClientOption func(*Client)

// WithBaseURL sets a custom base URL
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithTimeout sets a custom timeout
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.httpClient.Timeout = timeout
	}
}

// WithDebug enables debug mode to print HTTP request and response details
func WithDebug(debug bool) ClientOption {
	return func(c *Client) {
		c.debug = debug
	}
}

// hasSecret returns true if API Key Secret is provided
func (c *Client) hasSecret() bool {
	return c.apiKeySecret != ""
}

// getAuthHeader returns the Basic Auth header value
// Only used when API Key Secret is provided
func (c *Client) getAuthHeader() string {
	auth := c.apiKey + ":" + c.apiKeySecret
	return "Basic " + base64.StdEncoding.EncodeToString([]byte(auth))
}

// doRequest performs an HTTP request and returns the response
func (c *Client) doRequest(ctx context.Context, method, endpoint string, body io.Reader) (*http.Response, error) {
	url := c.baseURL + endpoint
	req, err := http.NewRequestWithContext(ctx, method, url, body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set Authorization header only if API Key Secret is provided (Basic Auth)
	// Otherwise, API Key will be included in the URL path
	if c.hasSecret() {
		req.Header.Set("Authorization", c.getAuthHeader())
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Debug: Print request details
	if c.debug {
		c.logRequest(req, body)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		if c.debug {
			log.Printf("[DEBUG] Request failed: %v\n", err)
		}
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}

	// Debug: Print response headers (body will be logged in doJSONRequest)
	if c.debug {
		c.logResponseHeaders(resp)
	}

	return resp, nil
}

// doJSONRequest performs a JSON request and unmarshals the response
func (c *Client) doJSONRequest(ctx context.Context, method, endpoint string, body interface{}, result interface{}) error {
	var bodyReader io.Reader
	var bodyBytes []byte
	if body != nil {
		var err error
		bodyBytes, err = json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	resp, err := c.doRequest(ctx, method, endpoint, bodyReader)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read response body for debug and error handling
	respBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// Debug: Print response body
	if c.debug {
		c.logResponseBody(respBodyBytes)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBodyBytes))
	}

	if result != nil {
		if err := json.Unmarshal(respBodyBytes, result); err != nil {
			if c.debug {
				log.Printf("[DEBUG] Failed to unmarshal response: %v\n", err)
			}
			return fmt.Errorf("failed to decode response: %w", err)
		}
		if c.debug {
			resultBytes, _ := json.MarshalIndent(result, "", "  ")
			log.Printf("[DEBUG] Parsed response object:\n%s\n", string(resultBytes))
		}
	}

	return nil
}

// logRequest logs detailed HTTP request information
func (c *Client) logRequest(req *http.Request, body io.Reader) {
	log.Printf("[DEBUG] ========== HTTP Request ==========\n")
	log.Printf("[DEBUG] Method: %s\n", req.Method)
	log.Printf("[DEBUG] URL: %s\n", req.URL.String())
	log.Printf("[DEBUG] Protocol: %s\n", req.Proto)
	log.Printf("[DEBUG] Host: %s\n", req.Host)

	log.Printf("[DEBUG] Headers:\n")
	for key, values := range req.Header {
		for _, value := range values {
			// Mask Authorization header for security
			if key == "Authorization" {
				log.Printf("[DEBUG]   %s: %s\n", key, maskAuthHeader(value))
			} else {
				log.Printf("[DEBUG]   %s: %s\n", key, value)
			}
		}
	}

	if body != nil {
		bodyBytes, err := io.ReadAll(body)
		if err == nil {
			// Create a new reader for the actual request since we consumed the body
			req.Body = io.NopCloser(bytes.NewReader(bodyBytes))

			var bodyStr string
			if len(bodyBytes) > 0 {
				var prettyJSON bytes.Buffer
				if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err == nil {
					bodyStr = prettyJSON.String()
				} else {
					bodyStr = string(bodyBytes)
				}
			}
			if bodyStr != "" {
				log.Printf("[DEBUG] Request Body:\n%s\n", bodyStr)
			}
		}
	}
	log.Printf("[DEBUG] ====================================\n")
}

// logResponseHeaders logs HTTP response headers
func (c *Client) logResponseHeaders(resp *http.Response) {
	log.Printf("[DEBUG] ========== HTTP Response Headers ==========\n")
	log.Printf("[DEBUG] Status: %s\n", resp.Status)
	log.Printf("[DEBUG] Status Code: %d\n", resp.StatusCode)
	log.Printf("[DEBUG] Protocol: %s\n", resp.Proto)

	log.Printf("[DEBUG] Headers:\n")
	for key, values := range resp.Header {
		for _, value := range values {
			log.Printf("[DEBUG]   %s: %s\n", key, value)
		}
	}
	log.Printf("[DEBUG] ============================================\n")
}

// logResponseBody logs HTTP response body
func (c *Client) logResponseBody(bodyBytes []byte) {
	log.Printf("[DEBUG] ========== HTTP Response Body ==========\n")
	if len(bodyBytes) > 0 {
		var prettyJSON bytes.Buffer
		if err := json.Indent(&prettyJSON, bodyBytes, "", "  "); err == nil {
			log.Printf("%s\n", prettyJSON.String())
		} else {
			log.Printf("%s\n", string(bodyBytes))
		}
	} else {
		log.Printf("[DEBUG] (empty body)\n")
	}
	log.Printf("[DEBUG] ===========================================\n")
}

// maskAuthHeader masks the authorization header for security
func maskAuthHeader(auth string) string {
	if len(auth) > 20 {
		return auth[:10] + "..." + auth[len(auth)-7:]
	}
	return "***"
}
