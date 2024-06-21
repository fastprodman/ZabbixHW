package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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
		name           string
		input          map[string]interface{}
		expectedCode   int
		expectedBody   map[string]interface{}
		expectedErrMsg string
	}{
		{
			name: "Valid JSON without id",
			input: map[string]interface{}{
				"name": "John",
			},
			expectedCode: http.StatusOK,
			expectedBody: map[string]interface{}{
				"name": "John",
				"id":   10,
			},
		},
		{
			name: "JSON with id",
			input: map[string]interface{}{
				"id":   1,
				"name": "John",
			},
			expectedCode:   http.StatusBadRequest,
			expectedErrMsg: "Field 'id' is not allowed\n",
		},
		{
			name:           "Invalid JSON",
			input:          nil,
			expectedCode:   http.StatusBadRequest,
			expectedErrMsg: "Error decoding JSON\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req *http.Request
			if tt.input != nil {
				body, _ := json.Marshal(tt.input)
				req = httptest.NewRequest("POST", "/record", bytes.NewBuffer(body))
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
			if tt.expectedBody != nil {

				var responseBody map[string]interface{}
				err := json.NewDecoder(rr.Body).Decode(&responseBody)
				if err != nil {
					t.Fatalf("error decoding response body: %v", err)
				}
				// Check if they have the same length
				if len(responseBody) != len(tt.expectedBody) {
					t.Error()
				}
			} else if tt.expectedErrMsg != "" {
				if rr.Body.String() != tt.expectedErrMsg {
					t.Errorf("expected error message %q, got %q", tt.expectedErrMsg, rr.Body.String())
				}
			}
		})
	}
}
