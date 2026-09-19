package db

import (
	"database/sql"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX idx_scheduler_date ON scheduler (date);
`

func Init(dbFIle string) error {
	_, err := os.Stat(dbFIle)
	var install bool
	if err != nil {
		install = true
	}

	var dbErr error
	DB, dbErr = sql.Open("sqlite", dbFIle)
	if dbErr != nil {
		return dbErr
	}

	if install {
		_, err = DB.Exec(schema)
		if err != nil {
			return err
		}
	}
	return nil
}
