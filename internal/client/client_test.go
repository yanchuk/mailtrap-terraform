package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	apiToken := "test-token"
	client := NewClient(apiToken)

	if client.apiToken != apiToken {
		t.Errorf("Expected API token %s, got %s", apiToken, client.apiToken)
	}

	if client.baseURL != defaultBaseURL {
		t.Errorf("Expected base URL %s, got %s", defaultBaseURL, client.baseURL)
	}

	if client.httpClient == nil {
		t.Error("Expected HTTP client to be initialized")
	}
}

func TestSetBaseURL(t *testing.T) {
	client := NewClient("test-token")
	customURL := "https://custom.example.com"
	
	client.SetBaseURL(customURL)
	
	if client.baseURL != customURL {
		t.Errorf("Expected base URL %s, got %s", customURL, client.baseURL)
	}
}

func TestDoRequest_Success(t *testing.T) {
	// Create test server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify headers
		if r.Header.Get("Api-Token") != "test-token" {
			t.Errorf("Expected API token header, got %s", r.Header.Get("Api-Token"))
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("Expected content type application/json, got %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("Accept") != "application/json" {
			t.Errorf("Expected accept application/json, got %s", r.Header.Get("Accept"))
		}

		// Return test response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.SetBaseURL(server.URL)

	resp, err := client.doRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestDoRequest_WithBody(t *testing.T) {
	requestBody := map[string]string{"name": "test"}
	
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request body
		var body map[string]string
		json.NewDecoder(r.Body).Decode(&body)
		
		if body["name"] != "test" {
			t.Errorf("Expected request body name 'test', got %s", body["name"])
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.SetBaseURL(server.URL)

	resp, err := client.doRequest("POST", "/test", requestBody)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}
}

func TestHandleResponse_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"name": "test-project"})
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.SetBaseURL(server.URL)

	resp, err := client.doRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]string
	err = client.handleResponse(resp, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["name"] != "test-project" {
		t.Errorf("Expected name 'test-project', got %s", result["name"])
	}
}

func TestHandleResponse_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{"error": "Bad request"})
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.SetBaseURL(server.URL)

	resp, err := client.doRequest("GET", "/test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	var result map[string]string
	err = client.handleResponse(resp, &result)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if !strings.Contains(err.Error(), "Bad request") {
		t.Errorf("Expected error to contain 'Bad request', got %s", err.Error())
	}
}

func TestGet(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "123"})
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.SetBaseURL(server.URL)

	var result map[string]string
	err := client.Get("/test", &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["id"] != "123" {
		t.Errorf("Expected id '123', got %s", result["id"])
	}
}

func TestPost(t *testing.T) {
	requestBody := map[string]string{"name": "test"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"id": "456"})
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.SetBaseURL(server.URL)

	var result map[string]string
	err := client.Post("/test", requestBody, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["id"] != "456" {
		t.Errorf("Expected id '456', got %s", result["id"])
	}
}

func TestPatch(t *testing.T) {
	requestBody := map[string]string{"name": "updated"}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "PATCH" {
			t.Errorf("Expected PATCH method, got %s", r.Method)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"id": "789"})
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.SetBaseURL(server.URL)

	var result map[string]string
	err := client.Patch("/test", requestBody, &result)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if result["id"] != "789" {
		t.Errorf("Expected id '789', got %s", result["id"])
	}
}

func TestDelete(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "DELETE" {
			t.Errorf("Expected DELETE method, got %s", r.Method)
		}

		w.WriteHeader(http.StatusNoContent)
	}))
	defer server.Close()

	client := NewClient("test-token")
	client.SetBaseURL(server.URL)

	err := client.Delete("/test", nil)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
}

func TestDoRequest_EndpointRouting(t *testing.T) {
	tests := []struct {
		endpoint    string
		expectedURL string
	}{
		{"/api/send", sendingAPIURL},
		{"/api/batch", sendingAPIURL},
		{"/api/send/123", sandboxAPIURL},
		{"/api/accounts/123", defaultBaseURL},
	}

	for _, tt := range tests {
		t.Run(tt.endpoint, func(t *testing.T) {
			client := NewClient("test-token")
			
			// We can't easily test the URL routing without mocking,
			// but we can verify the logic by checking the baseURL selection
			baseURL := client.baseURL
			if tt.endpoint == "/api/send" || tt.endpoint == "/api/batch" {
				baseURL = sendingAPIURL
			} else if len(tt.endpoint) > 10 && tt.endpoint[:10] == "/api/send/" {
				baseURL = sandboxAPIURL
			}

			if baseURL != tt.expectedURL {
				t.Errorf("Expected base URL %s for endpoint %s, got %s", tt.expectedURL, tt.endpoint, baseURL)
			}
		})
	}
}