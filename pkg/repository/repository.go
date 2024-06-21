package repository

type DatabaseRepo interface {
	CreateRecord()
	ReadRecord()
	UpdateRecord()
	DeleteRecord()
}
