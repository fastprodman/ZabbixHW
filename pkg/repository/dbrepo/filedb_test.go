package dbrepo

import (
	"encoding/json"
	"errors"
	"os"
	"reflect"
	"testing"
)

func Test_NewFileDB(t *testing.T) {
	t.Run("Successful DB creation", func(t *testing.T) {
		filepath := "./../../../testdata/readRecord.jsonl"
		file, err := os.OpenFile(filepath, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			t.Errorf("Failed to open file for test")
		}
		defer file.Close()

		// Test initializing the FileDB
		db, err := NewFileDB(file)
		if err != nil {
			t.Fatalf("Failed to initialize FileDB: %v", err)
		}

		// Check that the file and mutex are correctly initialized
		if db.file == nil {
			t.Errorf("FileDB.file is nil")
		}
		if db.fileMutex == nil {
			t.Errorf("FileDB.fileMutex is nil")
		}
	})
}

func Test_ReadRecord(t *testing.T) {
	tests := []struct {
		name           string
		recordID       uint32
		expectedRecord string
		expectedErr    error
	}{
		{
			name:           "Successful record read",
			recordID:       1,
			expectedRecord: `{"id": 1, "name": "Alice", "details": {"age": 30, "city": "New York"}}`,
			expectedErr:    nil,
		},
		{
			name:           "Record not existing",
			recordID:       999,
			expectedRecord: "",
			expectedErr:    ErrRecordNotFound,
		},
	}

	filepath := "./../../../testdata/readRecord.jsonl"
	file, err := os.OpenFile(filepath, os.O_RDONLY, 0)
	if err != nil {
		t.Errorf("Failed to open file for test")
	}
	defer file.Close()

	// Test initializing the FileDB
	db, err := NewFileDB(file)
	if err != nil {
		t.Fatalf("Failed to initialize FileDB: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			actualRecord, err := db.ReadRecord(tt.recordID)

			if !errors.Is(err, tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err)
			}

			if tt.expectedRecord != "" {
				var expectedBodyMap map[string]interface{}
				err = json.Unmarshal([]byte(tt.expectedRecord), &expectedBodyMap)
				if err != nil {
					t.Fatalf("error unmarshaling expected body: %v", err)
				}

				if !reflect.DeepEqual(expectedBodyMap, actualRecord) {
					t.Errorf("expected body %v, got %v", expectedBodyMap, actualRecord)
				}
			} else {
				if actualRecord != nil {
					t.Errorf("expected no record, but got %v", actualRecord)
				}
			}
		})
	}
}