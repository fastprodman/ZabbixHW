package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"zabbixhw/pkg/repository/dbrepo"
)

func Test_postRecordHandler(t *testing.T) {
	app := &application{
		DB: &dbrepo.TestDB{},
	}

	// Test cases
	tests := []struct {
		name           string
		input          string
		expectedCode   int
		expectedBody   string
		expectedErrMsg string
	}{
		{
			name: "Valid JSON without id",
			input: `{
				"name": "John"
			}`,
			expectedCode: http.StatusOK,
			expectedBody: `{
				"name": "John",
				"id": 10
			}`,
		},
		{
			name: "JSON with id",
			input: `{
				"id": 1,
				"name": "John"
			}`,
			expectedCode:   http.StatusBadRequest,
			expectedErrMsg: "Field 'id' is not allowed\n",
		},
		{
			name:           "No JSON given",
			input:          ``,
			expectedCode:   http.StatusBadRequest,
			expectedErrMsg: "Error decoding JSON\n",
		},
		{
			name: "Invalid JSON",
			input: `{
				"id": 1,
				"name
			}`,
			expectedCode:   http.StatusBadRequest,
			expectedErrMsg: "Error decoding JSON\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.input != "" {
				req = httptest.NewRequest("POST", "/record", bytes.NewBufferString(tt.input))
			} else {
				req = httptest.NewRequest("POST", "/record", nil)
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

			fmt.Println(rr.Body)

			// Check the response body
			if tt.expectedBody != "" {
				var expectedBodyMap map[string]interface{}
				err := json.Unmarshal([]byte(tt.expectedBody), &expectedBodyMap)
				if err != nil {
					t.Fatalf("error unmarshalling expected body: %v", err)
				}

				var responseBody map[string]interface{}
				err = json.NewDecoder(rr.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("error decoding response body: %v", err)
				}

				// Check if they have the same length
				if len(responseBody) != len(expectedBodyMap) {
					t.Errorf("expected response body length %d, got %d", len(expectedBodyMap), len(responseBody))
				}

				// Check if they are equal
				if !reflect.DeepEqual(responseBody, expectedBodyMap) {
					t.Errorf("expected response body %v, got %v", expectedBodyMap, responseBody)
				}
			} else if tt.expectedErrMsg != "" {
				if rr.Body.String() != tt.expectedErrMsg {
					t.Errorf("expected error message %q, got %q", tt.expectedErrMsg, rr.Body.String())
				}
			}
		})
	}
}
