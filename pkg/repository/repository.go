package repository

type DatabaseRepo interface {
	CreateRecord(map[string]interface{}) (map[string]interface{}, error)
	ReadRecord(id uint32) (map[string]interface{}, error)
	UpdateRecord()
	DeleteRecord()
}
