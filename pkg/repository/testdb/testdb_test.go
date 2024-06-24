package testdb

import (
	"reflect"
	"testing"
)

// TestCreateRecord tests the CreateRecord function with various scenarios
func Test_CreateRecord(t *testing.T) {
	t.Run("Create first record", func(t *testing.T) {
		db := &TestDB{}

		record := map[string]interface{}{
			"name": "Alice",
			"age":  30,
		}
		if err := db.CreateRecord(record); err != nil {
			t.Errorf("CreateRecord failed: %v", err)
		}
		if len(db.Data) != 1 {
			t.Errorf("expected 1 record, got %d", len(db.Data))
		}
		if id, ok := db.Data[0]["id"].(uint32); !ok || id != 1 {
			t.Errorf("expected id to be 1, got %v", db.Data[0]["id"])
		}

		expectedRecord := map[string]interface{}{
			"id":   uint32(1),
			"name": "Alice",
			"age":  30,
		}
		if !reflect.DeepEqual(db.Data[0], expectedRecord) {
			t.Errorf("expected record to be %v, got %v", expectedRecord, db.Data[0])
		}
	})

	t.Run("Create second record", func(t *testing.T) {
		db := &TestDB{
			Data: []map[string]interface{}{
				{
					"id":   uint32(1),
					"name": "Alice",
					"age":  30,
				},
			},
		}

		record := map[string]interface{}{
			"name": "Bob",
			"age":  25,
		}
		if err := db.CreateRecord(record); err != nil {
			t.Errorf("CreateRecord failed: %v", err)
		}
		if len(db.Data) != 2 {
			t.Errorf("expected 2 records, got %d", len(db.Data))
		}
		if id, ok := db.Data[1]["id"].(uint32); !ok || id != 2 {
			t.Errorf("expected id to be 2, got %v", db.Data[1]["id"])
		}

		expectedRecord := map[string]interface{}{
			"id":   uint32(2),
			"name": "Bob",
			"age":  25,
		}
		if !reflect.DeepEqual(db.Data[1], expectedRecord) {
			t.Errorf("expected record to be %v, got %v", expectedRecord, db.Data[1])
		}
	})

	t.Run("Create record with invalid ID in last record", func(t *testing.T) {
		db := &TestDB{
			Data: []map[string]interface{}{
				{
					"id":   "invalid", // invalid ID type
					"name": "Charlie",
					"age":  40,
				},
			},
		}

		record := map[string]interface{}{
			"name": "David",
			"age":  35,
		}
		err := db.CreateRecord(record)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if err.Error() != "invalid ID type in last record" {
			t.Errorf("expected error 'invalid ID type in last record', got %v", err)
		}
	})
}

func Test_ReadRecord(t *testing.T) {
	t.Run("Read existing record", func(t *testing.T) {
		db := &TestDB{
			Data: []map[string]interface{}{
				{
					"id":   uint32(1),
					"name": "Alice",
					"age":  30,
				},
				{
					"id":   uint32(2),
					"name": "Bob",
					"age":  25,
				},
			},
		}

		record, err := db.ReadRecord(1)
		if err != nil {
			t.Errorf("ReadRecord failed: %v", err)
		}

		expectedRecord := map[string]interface{}{
			"id":   uint32(1),
			"name": "Alice",
			"age":  30,
		}
		if !reflect.DeepEqual(record, expectedRecord) {
			t.Errorf("expected record to be %v, got %v", expectedRecord, record)
		}
	})

	t.Run("Read non-existing record", func(t *testing.T) {
		db := &TestDB{
			Data: []map[string]interface{}{
				{
					"id":   uint32(1),
					"name": "Alice",
					"age":  30,
				},
			},
		}

		_, err := db.ReadRecord(2)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if err.Error() != "record not found" {
			t.Errorf("expected error 'record not found', got %v", err)
		}
	})

	t.Run("Read record with invalid ID type", func(t *testing.T) {
		db := &TestDB{
			Data: []map[string]interface{}{
				{
					"id":   "invalid", // invalid ID type
					"name": "Charlie",
					"age":  40,
				},
			},
		}

		_, err := db.ReadRecord(1)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if err.Error() != "invalid ID type in record" {
			t.Errorf("expected error 'invalid ID type in record', got %v", err)
		}
	})
}

// TestUpdateRecord tests the UpdateRecord function with various scenarios
func Test_UpdateRecord(t *testing.T) {
	t.Run("Update existing record", func(t *testing.T) {
		db := &TestDB{
			Data: []map[string]interface{}{
				{
					"id":   uint32(1),
					"name": "Alice",
					"age":  30,
				},
				{
					"id":   uint32(2),
					"name": "Bob",
					"age":  25,
				},
			},
		}

		updatedData := map[string]interface{}{
			"name": "Alice Updated",
			"age":  31,
		}

		err := db.UpdateRecord(1, updatedData)
		if err != nil {
			t.Errorf("UpdateRecord failed: %v", err)
		}

		expectedRecord := map[string]interface{}{
			"id":   uint32(1),
			"name": "Alice Updated",
			"age":  31,
		}
		if !reflect.DeepEqual(db.Data[0], expectedRecord) {
			t.Errorf("expected record to be %v, got %v", expectedRecord, db.Data[0])
		}
	})

	t.Run("Update non-existing record", func(t *testing.T) {
		db := &TestDB{
			Data: []map[string]interface{}{
				{
					"id":   uint32(1),
					"name": "Alice",
					"age":  30,
				},
			},
		}

		updatedData := map[string]interface{}{
			"name": "Bob",
			"age":  25,
		}

		err := db.UpdateRecord(2, updatedData)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if err.Error() != "record not found" {
			t.Errorf("expected error 'record not found', got %v", err)
		}
	})

	t.Run("Update record with invalid ID type", func(t *testing.T) {
		db := &TestDB{
			Data: []map[string]interface{}{
				{
					"id":   "invalid", // invalid ID type
					"name": "Charlie",
					"age":  40,
				},
			},
		}

		updatedData := map[string]interface{}{
			"name": "Charlie Updated",
			"age":  41,
		}

		err := db.UpdateRecord(1, updatedData)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if err.Error() != "invalid ID type in record" {
			t.Errorf("expected error 'invalid ID type in record', got %v", err)
		}
	})
}

// TestDeleteRecord tests the DeleteRecord function with various scenarios
func TestDeleteRecord(t *testing.T) {
	t.Run("Delete existing record", func(t *testing.T) {
		db := &TestDB{
			Data: []map[string]interface{}{
				{
					"id":   uint32(1),
					"name": "Alice",
					"age":  30,
				},
				{
					"id":   uint32(2),
					"name": "Bob",
					"age":  25,
				},
			},
		}

		err := db.DeleteRecord(1)
		if err != nil {
			t.Errorf("DeleteRecord failed: %v", err)
		}

		if len(db.Data) != 1 {
			t.Errorf("expected 1 record, got %d", len(db.Data))
		}

		if db.Data[0]["id"] != uint32(2) {
			t.Errorf("expected remaining record ID to be 2, got %v", db.Data[0]["id"])
		}
	})

	t.Run("Delete non-existing record", func(t *testing.T) {
		db := &TestDB{
			Data: []map[string]interface{}{
				{
					"id":   uint32(1),
					"name": "Alice",
					"age":  30,
				},
			},
		}

		err := db.DeleteRecord(2)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if err.Error() != "record not found" {
			t.Errorf("expected error 'record not found', got %v", err)
		}
	})

	t.Run("Delete record with invalid ID type", func(t *testing.T) {
		db := &TestDB{
			Data: []map[string]interface{}{
				{
					"id":   "invalid", // invalid ID type
					"name": "Charlie",
					"age":  40,
				},
			},
		}

		err := db.DeleteRecord(1)
		if err == nil {
			t.Errorf("expected error, got nil")
		}
		if err.Error() != "invalid ID type in record" {
			t.Errorf("expected error 'invalid ID type in record', got %v", err)
		}
	})
}
