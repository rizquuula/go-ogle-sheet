package download_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"api/api/download"
)

func TestJsonHandler(t *testing.T) {
	// Define test cases
	tests := []struct {
		name           string
		queryParams    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Missing both parameters",
			queryParams:    "",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad Request: 'path' and 'gid' query parameter is required\n",
		},
		{
			name:           "Missing path parameter",
			queryParams:    "?gid=123",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad Request: 'path' and 'gid' query parameter is required\n",
		},
		{
			name:           "Missing gid parameter",
			queryParams:    "?path=test",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Bad Request: 'path' and 'gid' query parameter is required\n",
		},
		{
			name:           "All parameters provided",
			queryParams:    "?path=1_7lN9HpyUtmRwR28PGyrT6jZaHM3IdtCYcSdCscJkMQ&gid=0",
			expectedStatus: http.StatusOK,
			expectedBody:   `{"data":[{"Image":"https://example.com","Name":"Toyota","No":"1"},{"Image":"https://example.com/2","Name":"Daihatsu","No":"2"}]}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a request to pass to our handler with the specific query parameters
			req := httptest.NewRequest("GET", "/"+tt.queryParams, nil)
			// Create a ResponseRecorder to capture the response
			rr := httptest.NewRecorder()

			// Call the handler function directly
			download.JsonHandler(rr, req)

			// Check the status code
			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("expected status %v, got %v", tt.expectedStatus, status)
			}

			// Check the response body
			if body := rr.Body.String(); body != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, body)
			}
		})
	}
}