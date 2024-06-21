package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"zabbixhw/pkg/repository/dbrepo"
)

func TestRoutesExistence(t *testing.T) {
	app := &application{
		DB: &dbrepo.TestDB{},
	}
	handler := app.routes()

	tests := []struct {
		method string
		path   string
	}{
		{"POST", "/records"},
		{"GET", "/records/1"},
		{"PUT", "/records/1"},
		{"DELETE", "/records/1"},
	}

	for _, tt := range tests {
		t.Run(tt.method, func(t *testing.T) {
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			if rr.Code == http.StatusNotFound {
				t.Error("expected resource to exist, got 404")
			}
		})
	}
}
