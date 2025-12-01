package infura

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	client := NewClient("test-api-key", "test-api-secret")

	if client.apiKey != "test-api-key" {
		t.Errorf("Expected apiKey 'test-api-key', got '%s'", client.apiKey)
	}

	if client.apiKeySecret != "test-api-secret" {
		t.Errorf("Expected apiKeySecret 'test-api-secret', got '%s'", client.apiKeySecret)
	}

	if client.baseURL != BaseURL {
		t.Errorf("Expected baseURL '%s', got '%s'", BaseURL, client.baseURL)
	}

	if client.httpClient.Timeout != DefaultTimeout {
		t.Errorf("Expected timeout %v, got %v", DefaultTimeout, client.httpClient.Timeout)
	}
}

func TestNewClientWithOptions(t *testing.T) {
	customURL := "https://custom.url"
	customTimeout := 60 * time.Second

	client := NewClientWithOptions(
		"test-api-key",
		"test-api-secret",
		WithBaseURL(customURL),
		WithTimeout(customTimeout),
	)

	if client.baseURL != customURL {
		t.Errorf("Expected baseURL '%s', got '%s'", customURL, client.baseURL)
	}

	if client.httpClient.Timeout != customTimeout {
		t.Errorf("Expected timeout %v, got %v", customTimeout, client.httpClient.Timeout)
	}
}

func TestGetAuthHeader(t *testing.T) {
	client := NewClient("test-api-key", "test-api-secret")
	authHeader := client.getAuthHeader()

	if !strings.HasPrefix(authHeader, "Basic ") {
		t.Error("Auth header should start with 'Basic '")
	}

	// Decode and verify
	encoded := strings.TrimPrefix(authHeader, "Basic ")
	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		t.Fatalf("Failed to decode auth header: %v", err)
	}

	expected := "test-api-key:test-api-secret"
	if string(decoded) != expected {
		t.Errorf("Expected decoded auth '%s', got '%s'", expected, string(decoded))
	}
}

func TestDoRequest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Missing Authorization header")
		}

		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			t.Errorf("Expected Content-Type 'application/json', got '%s'", contentType)
		}

		accept := r.Header.Get("Accept")
		if accept != "application/json" {
			t.Errorf("Expected Accept 'application/json', got '%s'", accept)
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test": "response"}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("test-api-key", "test-api-secret", WithBaseURL(server.URL))

	resp, err := client.doRequest(context.Background(), "GET", "/test", nil)
	if err != nil {
		t.Fatalf("doRequest failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestDoJSONRequest(t *testing.T) {
	type TestResponse struct {
		Message string `json:"message"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("test-api-key", "test-api-secret", WithBaseURL(server.URL))

	var result TestResponse
	err := client.doJSONRequest(context.Background(), "GET", "/test", nil, &result)
	if err != nil {
		t.Fatalf("doJSONRequest failed: %v", err)
	}

	if result.Message != "success" {
		t.Errorf("Expected message 'success', got '%s'", result.Message)
	}
}

func TestDoJSONRequest_ErrorStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "bad request"}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("test-api-key", "test-api-secret", WithBaseURL(server.URL))

	var result map[string]interface{}
	err := client.doJSONRequest(context.Background(), "GET", "/test", nil, &result)
	if err == nil {
		t.Fatal("Expected error for bad request status but got nil")
	}

	if !strings.Contains(err.Error(), "400") {
		t.Errorf("Expected error to contain status code 400, got: %v", err)
	}
}

func TestNewClientWithAPIKey(t *testing.T) {
	client := NewClientWithAPIKey("test-api-key")

	if client.apiKey != "test-api-key" {
		t.Errorf("Expected apiKey 'test-api-key', got '%s'", client.apiKey)
	}

	if client.apiKeySecret != "" {
		t.Errorf("Expected empty apiKeySecret, got '%s'", client.apiKeySecret)
	}

	if client.hasSecret() {
		t.Error("hasSecret() should return false when secret is empty")
	}
}

func TestHasSecret(t *testing.T) {
	// Test with secret
	clientWithSecret := NewClient("test-api-key", "test-api-secret")
	if !clientWithSecret.hasSecret() {
		t.Error("hasSecret() should return true when secret is provided")
	}

	// Test without secret
	clientWithoutSecret := NewClientWithAPIKey("test-api-key")
	if clientWithoutSecret.hasSecret() {
		t.Error("hasSecret() should return false when secret is empty")
	}

	// Test with empty secret string
	clientEmptySecret := NewClient("test-api-key", "")
	if clientEmptySecret.hasSecret() {
		t.Error("hasSecret() should return false when secret is empty string")
	}
}

func TestDoRequest_APIKeyOnly(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify NO Authorization header when using API Key only
		auth := r.Header.Get("Authorization")
		if auth != "" {
			t.Error("Authorization header should not be present when using API Key only")
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"test": "response"}`))
	}))
	defer server.Close()

	client := NewClientWithAPIKeyAndOptions("test-api-key", WithBaseURL(server.URL))

	resp, err := client.doRequest(context.Background(), "GET", "/test", nil)
	if err != nil {
		t.Fatalf("doRequest failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestWithDebug(t *testing.T) {
	client := NewClientWithOptions("test-api-key", "test-api-secret", WithDebug(true))

	if !client.debug {
		t.Error("Expected debug to be true")
	}
}

func TestDoRequest_WithDebug(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Custom-Header", "test-value")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("test-api-key", "test-api-secret",
		WithBaseURL(server.URL),
		WithDebug(true))

	resp, err := client.doRequest(context.Background(), "GET", "/test", nil)
	if err != nil {
		t.Fatalf("doRequest failed: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode)
	}
}

func TestDoJSONRequest_WithDebug(t *testing.T) {
	type TestResponse struct {
		Message string `json:"message"`
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"message": "success"}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("test-api-key", "test-api-secret",
		WithBaseURL(server.URL),
		WithDebug(true))

	var result TestResponse
	err := client.doJSONRequest(context.Background(), "GET", "/test", nil, &result)
	if err != nil {
		t.Fatalf("doJSONRequest failed: %v", err)
	}

	if result.Message != "success" {
		t.Errorf("Expected message 'success', got '%s'", result.Message)
	}
}
