package repository

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
)

type Config struct {
	DBFileName string `yaml:"file"`
}

type Repository struct {
	db *sql.DB
}

func New(c Config) (*Repository, error) {
	r := &Repository{}
	db, err := sql.Open("sqlite3", c.DBFileName)
	if err != nil {
		return r, err
	}
	r.db = db
	err = r.checkDB()
	return r, err
}

func (r *Repository) checkDB() error {
	if r.db == nil {
		return errors.New("nil db pointer")
	}
	rows, err := r.db.Query("SELECT id, label FROM processes")
	if err != nil {
		return err
	}
	colsTypes, err := rows.ColumnTypes()
	if err != nil {
		return err
	}
	idColExists := false
	for _, ct := range colsTypes {
		if ct.Name() == "id" {
			idColExists = true
			if ct.DatabaseTypeName() != "INTEGER" {
				return errors.New("id column must be integer")
			}
		}
	}
	if !idColExists {
		return errors.New("id column not exists")
	}
	if err = rows.Close(); err != nil {
		return err
	}
	return nil
}
