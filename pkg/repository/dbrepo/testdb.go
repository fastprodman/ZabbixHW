package dbrepo

type TestDB struct {
	data []map[string]interface{}
}

func (db *TestDB) CreateRecord(data map[string]interface{}) (map[string]interface{}, error) {
	data["id"] = 10
	return data, nil
}

func (db *TestDB) ReadRecord()   {}
func (db *TestDB) UpdateRecord() {}
func (db *TestDB) DeleteRecord() {}
