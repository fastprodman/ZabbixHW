package repository

type DatabaseRepo interface {
	CreateRecord(map[string]interface{}) (map[string]interface{}, error)
	ReadRecord()
	UpdateRecord()
	DeleteRecord()
}
