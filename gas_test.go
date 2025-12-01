package infura

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestGetSuggestedGasFees(t *testing.T) {
	// Mock response data
	mockResponse := SuggestedGasFees{
		Low: GasFeeLevel{
			SuggestedMaxPriorityFeePerGas: "0.05",
			SuggestedMaxFeePerGas:         "24.086058416",
			MinWaitTimeEstimate:           15000,
			MaxWaitTimeEstimate:           30000,
		},
		Medium: GasFeeLevel{
			SuggestedMaxPriorityFeePerGas: "0.1",
			SuggestedMaxFeePerGas:         "32.548678862",
			MinWaitTimeEstimate:           15000,
			MaxWaitTimeEstimate:           45000,
		},
		High: GasFeeLevel{
			SuggestedMaxPriorityFeePerGas: "0.3",
			SuggestedMaxFeePerGas:         "41.161299308",
			MinWaitTimeEstimate:           15000,
			MaxWaitTimeEstimate:           60000,
		},
		EstimatedBaseFee:           "24.036058416",
		NetworkCongestion:          0.7143,
		LatestPriorityFeeRange:     []string{"0.1", "20"},
		HistoricalPriorityFeeRange: []string{"0.007150439", "113"},
		HistoricalBaseFeeRange:     []string{"19.531410688", "36.299069766"},
		PriorityFeeTrend:           "down",
		BaseFeeTrend:               "down",
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		expectedPath := "/networks/1/suggestedGasFees"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify Basic Auth header
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Missing Authorization header")
		}

		// Set response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewClientWithOptions("test-api-key", "test-api-secret", WithBaseURL(server.URL))

	// Test GetSuggestedGasFees
	result, err := client.GetSuggestedGasFees(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetSuggestedGasFees failed: %v", err)
	}

	// Verify response
	if result.Low.SuggestedMaxPriorityFeePerGas != mockResponse.Low.SuggestedMaxPriorityFeePerGas {
		t.Errorf("Expected Low.SuggestedMaxPriorityFeePerGas %s, got %s",
			mockResponse.Low.SuggestedMaxPriorityFeePerGas,
			result.Low.SuggestedMaxPriorityFeePerGas)
	}

	if result.Medium.SuggestedMaxFeePerGas != mockResponse.Medium.SuggestedMaxFeePerGas {
		t.Errorf("Expected Medium.SuggestedMaxFeePerGas %s, got %s",
			mockResponse.Medium.SuggestedMaxFeePerGas,
			result.Medium.SuggestedMaxFeePerGas)
	}

	if result.High.MinWaitTimeEstimate != mockResponse.High.MinWaitTimeEstimate {
		t.Errorf("Expected High.MinWaitTimeEstimate %d, got %d",
			mockResponse.High.MinWaitTimeEstimate,
			result.High.MinWaitTimeEstimate)
	}

	if result.EstimatedBaseFee != mockResponse.EstimatedBaseFee {
		t.Errorf("Expected EstimatedBaseFee %s, got %s",
			mockResponse.EstimatedBaseFee,
			result.EstimatedBaseFee)
	}

	if result.NetworkCongestion != mockResponse.NetworkCongestion {
		t.Errorf("Expected NetworkCongestion %f, got %f",
			mockResponse.NetworkCongestion,
			result.NetworkCongestion)
	}
}

func TestGetSuggestedGasFees_ErrorResponse(t *testing.T) {
	// Create mock server that returns an error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Unauthorized"}`))
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewClientWithOptions("invalid-key", "invalid-secret", WithBaseURL(server.URL))

	// Test GetSuggestedGasFees with error response
	_, err := client.GetSuggestedGasFees(context.Background(), 1)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
}

func TestGetSuggestedGasFees_InvalidJSON(t *testing.T) {
	// Create mock server that returns invalid JSON
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`invalid json`))
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewClientWithOptions("test-api-key", "test-api-secret", WithBaseURL(server.URL))

	// Test GetSuggestedGasFees with invalid JSON
	_, err := client.GetSuggestedGasFees(context.Background(), 1)
	if err == nil {
		t.Fatal("Expected error for invalid JSON but got nil")
	}
}

func TestGetSuggestedGasFees_APIKeyOnly(t *testing.T) {
	// Mock response data
	mockResponse := SuggestedGasFees{
		Low: GasFeeLevel{
			SuggestedMaxPriorityFeePerGas: "0.05",
			SuggestedMaxFeePerGas:         "24.086058416",
			MinWaitTimeEstimate:           15000,
			MaxWaitTimeEstimate:           30000,
		},
		Medium: GasFeeLevel{
			SuggestedMaxPriorityFeePerGas: "0.1",
			SuggestedMaxFeePerGas:         "32.548678862",
			MinWaitTimeEstimate:           15000,
			MaxWaitTimeEstimate:           45000,
		},
		High: GasFeeLevel{
			SuggestedMaxPriorityFeePerGas: "0.3",
			SuggestedMaxFeePerGas:         "41.161299308",
			MinWaitTimeEstimate:           15000,
			MaxWaitTimeEstimate:           60000,
		},
		EstimatedBaseFee:           "24.036058416",
		NetworkCongestion:          0.7143,
		LatestPriorityFeeRange:     []string{"0.1", "20"},
		HistoricalPriorityFeeRange: []string{"0.007150439", "113"},
		HistoricalBaseFeeRange:     []string{"19.531410688", "36.299069766"},
		PriorityFeeTrend:           "down",
		BaseFeeTrend:               "down",
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		// When using API Key only, the path should include /v3/{apiKey}/
		expectedPath := "/v3/test-api-key/networks/1/suggestedGasFees"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify NO Authorization header when using API Key only
		auth := r.Header.Get("Authorization")
		if auth != "" {
			t.Error("Authorization header should not be present when using API Key only")
		}

		// Set response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with only API Key (no secret)
	client := NewClientWithAPIKeyAndOptions("test-api-key", WithBaseURL(server.URL))

	// Test GetSuggestedGasFees
	result, err := client.GetSuggestedGasFees(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetSuggestedGasFees failed: %v", err)
	}

	// Verify response
	if result.Low.SuggestedMaxPriorityFeePerGas != mockResponse.Low.SuggestedMaxPriorityFeePerGas {
		t.Errorf("Expected Low.SuggestedMaxPriorityFeePerGas %s, got %s",
			mockResponse.Low.SuggestedMaxPriorityFeePerGas,
			result.Low.SuggestedMaxPriorityFeePerGas)
	}

	if result.Medium.SuggestedMaxFeePerGas != mockResponse.Medium.SuggestedMaxFeePerGas {
		t.Errorf("Expected Medium.SuggestedMaxFeePerGas %s, got %s",
			mockResponse.Medium.SuggestedMaxFeePerGas,
			result.Medium.SuggestedMaxFeePerGas)
	}
}

func TestGetSuggestedGasFees_APIKeyOnly_EmptySecret(t *testing.T) {
	// Mock response data
	mockResponse := SuggestedGasFees{
		Low: GasFeeLevel{
			SuggestedMaxPriorityFeePerGas: "0.05",
			SuggestedMaxFeePerGas:         "24.086058416",
			MinWaitTimeEstimate:           15000,
			MaxWaitTimeEstimate:           30000,
		},
		Medium: GasFeeLevel{
			SuggestedMaxPriorityFeePerGas: "0.1",
			SuggestedMaxFeePerGas:         "32.548678862",
			MinWaitTimeEstimate:           15000,
			MaxWaitTimeEstimate:           45000,
		},
		High: GasFeeLevel{
			SuggestedMaxPriorityFeePerGas: "0.3",
			SuggestedMaxFeePerGas:         "41.161299308",
			MinWaitTimeEstimate:           15000,
			MaxWaitTimeEstimate:           60000,
		},
		EstimatedBaseFee:           "24.036058416",
		NetworkCongestion:          0.7143,
		LatestPriorityFeeRange:     []string{"0.1", "20"},
		HistoricalPriorityFeeRange: []string{"0.007150439", "113"},
		HistoricalBaseFeeRange:     []string{"19.531410688", "36.299069766"},
		PriorityFeeTrend:           "down",
		BaseFeeTrend:               "down",
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// When secret is empty, should use URL path auth
		expectedPath := "/v3/test-api-key/networks/1/suggestedGasFees"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify NO Authorization header when secret is empty
		auth := r.Header.Get("Authorization")
		if auth != "" {
			t.Error("Authorization header should not be present when secret is empty")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with empty secret (should use API Key only method)
	client := NewClientWithOptions("test-api-key", "", WithBaseURL(server.URL))

	// Test GetSuggestedGasFees
	result, err := client.GetSuggestedGasFees(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetSuggestedGasFees failed: %v", err)
	}

	// Verify response
	if result.Low.SuggestedMaxPriorityFeePerGas != mockResponse.Low.SuggestedMaxPriorityFeePerGas {
		t.Errorf("Expected Low.SuggestedMaxPriorityFeePerGas %s, got %s",
			mockResponse.Low.SuggestedMaxPriorityFeePerGas,
			result.Low.SuggestedMaxPriorityFeePerGas)
	}
}

func TestGetBaseFeeHistory(t *testing.T) {
	// Mock response data - API returns array directly
	mockResponse := BaseFeeHistory{"24.036058416", "25.123456789", "23.987654321"}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify the request
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		expectedPath := "/networks/1/baseFeeHistory"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify Basic Auth header
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Missing Authorization header")
		}

		// Set response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewClientWithOptions("test-api-key", "test-api-secret", WithBaseURL(server.URL))

	// Test GetBaseFeeHistory
	result, err := client.GetBaseFeeHistory(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetBaseFeeHistory failed: %v", err)
	}

	// Verify response
	if len(result) != len(mockResponse) {
		t.Errorf("Expected BaseFeeHistory length %d, got %d",
			len(mockResponse),
			len(result))
	}

	if result[0] != mockResponse[0] {
		t.Errorf("Expected BaseFeeHistory[0] %s, got %s",
			mockResponse[0],
			result[0])
	}
}

func TestGetBaseFeeHistory_APIKeyOnly(t *testing.T) {
	// Mock response data - API returns array directly
	mockResponse := BaseFeeHistory{"24.036058416", "25.123456789"}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/v3/test-api-key/networks/1/baseFeeHistory"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify NO Authorization header when using API Key only
		auth := r.Header.Get("Authorization")
		if auth != "" {
			t.Error("Authorization header should not be present when using API Key only")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with only API Key
	client := NewClientWithAPIKeyAndOptions("test-api-key", WithBaseURL(server.URL))

	// Test GetBaseFeeHistory
	result, err := client.GetBaseFeeHistory(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetBaseFeeHistory failed: %v", err)
	}

	if len(result) != 2 {
		t.Errorf("Expected BaseFeeHistory length 2, got %d", len(result))
	}
}

func TestGetBaseFeePercentile(t *testing.T) {
	// Mock response data
	mockResponse := BaseFeePercentile{
		BaseFeePercentile: "50",
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		expectedPath := "/networks/1/baseFeePercentile"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify Basic Auth header
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Missing Authorization header")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewClientWithOptions("test-api-key", "test-api-secret", WithBaseURL(server.URL))

	// Test GetBaseFeePercentile
	result, err := client.GetBaseFeePercentile(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetBaseFeePercentile failed: %v", err)
	}

	// Verify response
	if result.BaseFeePercentile != mockResponse.BaseFeePercentile {
		t.Errorf("Expected BaseFeePercentile %s, got %s",
			mockResponse.BaseFeePercentile,
			result.BaseFeePercentile)
	}
}

func TestGetBaseFeePercentile_APIKeyOnly(t *testing.T) {
	// Mock response data
	mockResponse := BaseFeePercentile{
		BaseFeePercentile: "75",
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/v3/test-api-key/networks/1/baseFeePercentile"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		auth := r.Header.Get("Authorization")
		if auth != "" {
			t.Error("Authorization header should not be present when using API Key only")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with only API Key
	client := NewClientWithAPIKeyAndOptions("test-api-key", WithBaseURL(server.URL))

	// Test GetBaseFeePercentile
	result, err := client.GetBaseFeePercentile(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetBaseFeePercentile failed: %v", err)
	}

	if result.BaseFeePercentile != "75" {
		t.Errorf("Expected BaseFeePercentile '75', got %s", result.BaseFeePercentile)
	}
}

func TestGetBusyThreshold(t *testing.T) {
	// Mock response data
	mockResponse := BusyThreshold{
		BusyThreshold: "0.7",
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET request, got %s", r.Method)
		}

		expectedPath := "/networks/1/busyThreshold"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		// Verify Basic Auth header
		auth := r.Header.Get("Authorization")
		if auth == "" {
			t.Error("Missing Authorization header")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with mock server URL
	client := NewClientWithOptions("test-api-key", "test-api-secret", WithBaseURL(server.URL))

	// Test GetBusyThreshold
	result, err := client.GetBusyThreshold(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetBusyThreshold failed: %v", err)
	}

	// Verify response
	if result.BusyThreshold != mockResponse.BusyThreshold {
		t.Errorf("Expected BusyThreshold %s, got %s",
			mockResponse.BusyThreshold,
			result.BusyThreshold)
	}
}

func TestGetBusyThreshold_APIKeyOnly(t *testing.T) {
	// Mock response data
	mockResponse := BusyThreshold{
		BusyThreshold: "0.8",
	}

	// Create mock server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expectedPath := "/v3/test-api-key/networks/1/busyThreshold"
		if r.URL.Path != expectedPath {
			t.Errorf("Expected path %s, got %s", expectedPath, r.URL.Path)
		}

		auth := r.Header.Get("Authorization")
		if auth != "" {
			t.Error("Authorization header should not be present when using API Key only")
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(mockResponse)
	}))
	defer server.Close()

	// Create client with only API Key
	client := NewClientWithAPIKeyAndOptions("test-api-key", WithBaseURL(server.URL))

	// Test GetBusyThreshold
	result, err := client.GetBusyThreshold(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetBusyThreshold failed: %v", err)
	}

	if result.BusyThreshold != "0.8" {
		t.Errorf("Expected BusyThreshold '0.8', got %s", result.BusyThreshold)
	}
}

func TestGetBaseFeeHistory_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte(`{"error": "Unauthorized"}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("invalid-key", "invalid-secret", WithBaseURL(server.URL))

	_, err := client.GetBaseFeeHistory(context.Background(), 1)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
}

func TestGetBaseFeePercentile_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(`{"error": "Bad Request"}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("invalid-key", "invalid-secret", WithBaseURL(server.URL))

	_, err := client.GetBaseFeePercentile(context.Background(), 1)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
}

func TestGetBusyThreshold_ErrorResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(`{"error": "Not Found"}`))
	}))
	defer server.Close()

	client := NewClientWithOptions("invalid-key", "invalid-secret", WithBaseURL(server.URL))

	_, err := client.GetBusyThreshold(context.Background(), 1)
	if err == nil {
		t.Fatal("Expected error but got nil")
	}
}

func TestClient_GetSuggestedGasFees(t *testing.T) {
	client := NewClientWithAPIKey(os.Getenv("InfuraAPIKey"))
	data, err := client.GetSuggestedGasFees(context.Background(), 1)
	if err != nil {
		t.Fatalf("Error fetching suggested gas fees: %v", err)
	}
	dataStr, _ := json.MarshalIndent(data, "", "  ")
	t.Logf("Suggested Gas Fees: %+v", string(dataStr))
}

func TestClient_GetBaseFeeHistory(t *testing.T) {
	client := NewClientWithAPIKey(os.Getenv("InfuraAPIKey"))
	data, err := client.GetBaseFeeHistory(context.Background(), 1)
	if err != nil {
		t.Fatalf("Error fetching base fee history: %v", err)
	}
	dataStr, _ := json.MarshalIndent(data, "", "  ")
	t.Logf("Base Fee History: %+v", string(dataStr))
}

func TestClient_GetBaseFeePercentile(t *testing.T) {
	client := NewClientWithAPIKey(os.Getenv("InfuraAPIKey"))
	data, err := client.GetBaseFeePercentile(context.Background(), 1)
	if err != nil {
		t.Fatalf("Error fetching base fee percentile: %v", err)
	}
	dataStr, _ := json.MarshalIndent(data, "", "  ")
	t.Logf("Base Fee Percentile: %+v", string(dataStr))
}

func TestClient_GetBusyThreshold(t *testing.T) {
	client := NewClientWithAPIKey(os.Getenv("InfuraAPIKey"))
	data, err := client.GetBusyThreshold(context.Background(), 1)
	if err != nil {
		t.Fatalf("Error fetching busy threshold: %v", err)
	}
	dataStr, _ := json.MarshalIndent(data, "", "  ")
	t.Logf("Busy Threshold: %+v", string(dataStr))
}
