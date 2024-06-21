package dbrepo

type FileDB struct {
}

func (f *FileDB) CreateRecord(map[string]interface{}) (map[string]interface{}, error) {
	var m map[string]interface{}
	return m, nil
}

func (f *FileDB) ReadRecord()   {}
func (f *FileDB) UpdateRecord() {}
func (f *FileDB) DeleteRecord() {}
