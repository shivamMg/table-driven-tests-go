package api

type Database interface {
	CreateTODO(todo *TODO) error
}

type DBClient struct{}

func (db *DBClient) CreateTODO(todo *TODO) error {
	return nil
}
