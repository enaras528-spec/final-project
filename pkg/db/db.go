package db

import (
	"database/sql"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX IF NOT EXISTS idx_scheduler_date ON scheduler (date);
`

func Init(dbFile string) error {
	var err error
	DB, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	_, err = DB.Exec(schema)
	if err != nil {
		DB.Close()
		return err
	}

	return nil
}
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
