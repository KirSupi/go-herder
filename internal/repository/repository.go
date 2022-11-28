package repository

import (
	_ "github.com/mattn/go-sqlite3"
)

type Repository interface {
	IterProcesses() chan ProcessData
}

type ProcessData struct {
	ID      int
	Label   *string
	Command string
	Params  string
}
