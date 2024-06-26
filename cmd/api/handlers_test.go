package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"zabbixhw/pkg/helpers"
	"zabbixhw/pkg/repository/testdb"
)

func Test_postRecordHandler(t *testing.T) {
	// Test cases
	tests := []struct {
		name         string
		input        string
		expectedCode int
		expectedBody string
		errExpected  bool
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
			name: "Valid JSON nested structure",
			input: `{
				"name": "John",
				"address": {
					"street": "123 Main St",
					"city": "Anytown"
				}
			}`,
			expectedCode: http.StatusOK,
			expectedBody: `{"id":1,"name":"John","address":{"street":"123 Main St","city":"Anytown"}}`,
			errExpected:  false,
		},
		{
			name: "JSON with id",
			input: `{
				"id": 1,
				"name": "John"
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Field 'id' is not allowed\n",
			errExpected:  true,
		},
		{
			name:         "No JSON given",
			input:        ``,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Error decoding JSON\n",
			errExpected:  true,
		},
		{
			name: "Invalid JSON",
			input: `{
				"id": 1,
				"name
			}`,
			expectedCode: http.StatusBadRequest,
			expectedBody: "Error decoding JSON\n",
			errExpected:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &application{
				DB: &testdb.TestDB{
					Data: []map[string]interface{}{},
				},
			}

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

			if !tt.errExpected {
				var actualBodyMap map[string]interface{}

				err := json.Unmarshal(rr.Body.Bytes(), &actualBodyMap)
				if err != nil {
					t.Fatalf("error unmarshaling actual body: %v", err)
				}

				ok, err := helpers.CompareJSONWithMap(tt.expectedBody, actualBodyMap)
				if err != nil {
					t.Fatalf("error unmarshaling actual body: %v", err)
				}

				if !ok {
					t.Errorf("expected body %v, got %v", tt.expectedBody, actualBodyMap)
				}

			} else {
				// Check the response body
				if rr.Body.String() != tt.expectedBody {
					t.Errorf("expected body %q, got %q", tt.expectedBody, rr.Body.String())
				}
			}

		})
	}
}

func Test_getRecordHandler(t *testing.T) {
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
			app := &application{
				DB: &testdb.TestDB{
					Data: []map[string]interface{}{
						{"id": uint32(1), "name": "Record 1"},
					},
				},
			}
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
			app := &application{
				DB: &testdb.TestDB{
					Data: []map[string]interface{}{
						{"id": uint32(1), "name": "Record 1"},
					},
				},
			}
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
			app := &application{
				DB: &testdb.TestDB{
					Data: []map[string]interface{}{
						{"id": uint32(1), "name": "Record 1"},
					},
				},
			}
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
