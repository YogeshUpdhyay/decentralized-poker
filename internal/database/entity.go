package database

type Entity interface {
	GetByID(id string) (interface{}, error)
	Save() error
	Update() error
	Delete() error
}
