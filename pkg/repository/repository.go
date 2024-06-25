package repository

type DatabaseRepo interface {
	CreateRecord(data map[string]interface{}) error
	ReadRecord(id uint32) (map[string]interface{}, error)
	UpdateRecord(id uint32, data map[string]interface{}) error
	DeleteRecord(id uint32) error
}
