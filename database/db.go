package database

import (
	"database/sql"
	_ "github.com/lib/pq"
	)

type Database struct {
	config *Config
	connection *sql.DB
}

func New(config *Config) *Database{
	return &Database{
		config: config,
	}
}

func (database *Database) Open() error {
	db, err := sql.Open("postgres", database.config.DatabaseUrl)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	database.connection = db

	return nil
}

func (database *Database) Close() {
	database.connection.Close()
}



