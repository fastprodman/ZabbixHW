package dbrepo

import (
	"bufio"
	"encoding/json"
	"errors"
	"os"
	"testing"
	"zabbixhw/pkg/helpers"
)

func Test_NewFileDB(t *testing.T) {
	t.Run("Successful DB creation", func(t *testing.T) {
		filepath := "./../../../testdata/testBase.jsonl"
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

	filepath := "./../../../testdata/testBase.jsonl"
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
				equal, err := helpers.CompareJSONWithMap(tt.expectedRecord, actualRecord)
				if err != nil {
					t.Fatalf("got error comparing records: %s", err.Error())
				}

				if !equal {
					t.Errorf("records not equal as expected %v and %s", actualRecord, tt.expectedRecord)
				}
			} else {
				if actualRecord != nil {
					t.Errorf("expected no record, but got %v", actualRecord)
				}
			}
		})
	}
}

func Test_CreateRecord(t *testing.T) {
	t.Run("Successful record creation", func(t *testing.T) {
		// Create a temporary file for testing
		file, err := os.CreateTemp("./../../../testdata/", "testdb*.jsonl")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(file.Name())
		defer file.Close()

		// Initialize FileDB with the temp file
		db, err := NewFileDB(file)
		if err != nil {
			t.Fatalf("Failed to initialize FileDB: %v", err)
		}
		testRecord := `{"name":"Bob","details":{"age":25,"city":"Los Angeles"}}`
		expectedDBRecord := `{"id":1,"name":"Bob","details":{"age":25,"city":"Los Angeles"}}`

		var testRecordMap map[string]interface{}
		err = json.Unmarshal([]byte(testRecord), &testRecordMap)
		if err != nil {
			t.Fatalf("error unmarshaling expected body: %v", err)
		}

		err = db.CreateRecord(testRecordMap)
		if err != nil {
			t.Fatalf("Failed to create record: %v", err)
		}
		// Check if record modified as expected
		equals, err := helpers.CompareJSONWithMap(expectedDBRecord, testRecordMap)
		if err != nil {
			t.Fatalf("Failed to create record: %v", err)
		}

		if !equals {
			t.Errorf("Record not assigned with expected id")
		}

		// Check if record is in fle
		found, err := jsonRecordInFile(expectedDBRecord, file)
		if err != nil {
			t.Fatalf("Failed to create record: %v", err)
		}

		if !found {
			t.Errorf("Expected record not in db")
		}
	})
}

func jsonRecordInFile(record string, file *os.File) (bool, error) {
	file.Seek(0, 0)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		recordInFile := scanner.Text()

		found, err := helpers.CompareJSONStrings(recordInFile, record)
		if err != nil {
			return false, err
		}

		if found {
			return true, nil
		}

	}
	return false, nil
}
