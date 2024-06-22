package main

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
	"zabbixhw/pkg/repository/dbrepo"
)

func Test_postRecordHandler(t *testing.T) {
	app := &application{
		DB: &dbrepo.TestDB{},
	}

	// Test cases
	tests := []struct {
		name         string
		input        string
		expectedCode int
		expectedBody string
	}{
		{
			name: "Valid JSON without id",
			input: `{
				"name": "John"
			}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":1,"name":"John"}`,
		},
		{
			name: "JSON with id",
			input: `{
				"id": 1,
				"name": "John"
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Field 'id' is not allowed\n",
		},
		{
			name:         "No JSON given",
			input:        ``,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Error decoding JSON\n",
		},
		{
			name: "Invalid JSON",
			input: `{
				"id": 1,
				"name
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Error decoding JSON\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.input != "" {
				req = httptest.NewRequest("POST", "/records", bytes.NewBufferString(tt.input))
			} else {
				req = httptest.NewRequest("POST", "/records", nil)
			}

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Call the handler
			handler := http.HandlerFunc(app.postRecordHandler)
			handler.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, rr.Code)
			}

			// Check the response body
			if rr.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rr.Body.String())
			}

		})
	}
}

func Test_getRecordHandler(t *testing.T) {
	app := &application{
		DB: &dbrepo.TestDB{
			Data: []map[string]interface{}{
				{"id": uint32(1), "name": "Record 1"},
			},
		},
	}

	// Test cases
	tests := []struct {
		name         string
		path         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid ID",
			path:         "/records/1",
			expectedCode: http.StatusOK,
			expectedBody: `{"id":1,"name":"Record 1"}`,
		},
		{
			name:         "Invalid ID",
			path:         "/records/abc",
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid ID\n",
		},
		{
			name:         "Record Not Found",
			path:         "/records/999",
			expectedCode: http.StatusBadRequest,
			expectedBody: "record not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			rr := httptest.NewRecorder()
			mux := http.NewServeMux()
			mux.HandleFunc("GET /records/{id}", app.getRecordHandler)
			mux.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, rr.Code)
			}

			// Check the response body
			if rr.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func Test_deleteRecordHandler(t *testing.T) {
	app := &application{
		DB: &dbrepo.TestDB{
			Data: []map[string]interface{}{
				{"id": uint32(1), "name": "Record 1"},
			},
		},
	}

	// Test cases
	tests := []struct {
		name         string
		path         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Valid ID",
			path:         "/records/1",
			expectedCode: http.StatusNoContent,
			expectedBody: "",
		},
		{
			name:         "Invalid ID",
			path:         "/records/abc",
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid ID\n",
		},
		{
			name:         "Record Not Found",
			path:         "/records/999",
			expectedCode: http.StatusBadRequest,
			expectedBody: "record not found\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("DELETE", tt.path, nil)
			rr := httptest.NewRecorder()
			mux := http.NewServeMux()
			mux.HandleFunc("DELETE /records/{id}", app.deleteRecordHandler)
			mux.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, rr.Code)
			}

			// Check the response body
			if rr.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}

func Test_updateRecordHandler(t *testing.T) {
	app := &application{
		DB: &dbrepo.TestDB{
			Data: []map[string]interface{}{
				{"id": uint32(1), "name": "Record 1"},
			},
		},
	}

	tests := []struct {
		name         string
		path         string
		inputJSON    string
		expectedCode int
		expectedBody string
	}{
		{
			name: "Valid update",
			path: "/records/1",
			inputJSON: `{
				"name": "John",
				"pet": "dog"
			}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":1,"name":"John","pet":"dog"}`,
		},
		{
			name: "No record",
			path: "/records/2",
			inputJSON: `{
				"name": "John",
				"pet": "dog"
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "record not found\n",
		},
		{
			name: "Invalid Id",
			path: "/records/abc",
			inputJSON: `{
				"name": "John",
				"pet": "dog"
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Invalid ID\n",
		},
		{
			name: "Field Id in body",
			path: "/records/1",
			inputJSON: `{
				"id": 1,
				"name": "John",
				"pet": "dog"
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Field 'id' is not allowed\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.inputJSON != "" {
				req = httptest.NewRequest("PUT", tt.path, bytes.NewBufferString(tt.inputJSON))
			} else {
				req = httptest.NewRequest("PUT", tt.path, nil)
			}
			rr := httptest.NewRecorder()
			mux := http.NewServeMux()
			mux.HandleFunc("PUT /records/{id}", app.putRecordHandler)
			mux.ServeHTTP(rr, req)

			// Check the status code
			if rr.Code != tt.expectedCode {
				t.Errorf("expected status code %d, got %d", tt.expectedCode, rr.Code)
			}

			// Check the response body
			if rr.Body.String() != tt.expectedBody {
				t.Errorf("expected body %q, got %q", tt.expectedBody, rr.Body.String())
			}
		})
	}
}
