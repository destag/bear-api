package database

import "github.com/destag/bear-api/model"

func (db *Database) ListBears() ([]*model.Bear, error) {
	var bears []*model.Bear
	err := db.conn.Select(&bears, "SELECT * FROM bears")
	return bears, err
}

func (db *Database) CreateBear(name string) (*model.Bear, error) {
	res, err := db.conn.Exec("INSERT INTO bears (name) VALUES (?)", name)
	if err != nil {
		return nil, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.Bear{
		ID:   uint(id),
		Name: name,
	}, nil
}
