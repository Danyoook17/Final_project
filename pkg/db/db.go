package db

import (
	"database/sql"
	"strings"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE IF NOT EXISTS scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(255) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);

CREATE INDEX IF NOT EXISTS scheduler_date_idx ON scheduler(date);
`

func Init(dbFile string) error {
	d, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return err
	}

	if err := d.Ping(); err != nil {
		_ = d.Close()
		return err
	}

	parts := strings.Split(schema, ";")
	for _, p := range parts {
		stmt := strings.TrimSpace(p)
		if stmt == "" {
			continue
		}
		if _, err := d.Exec(stmt); err != nil {
			_ = d.Close()
			return err
		}
	}

	DB = d
	return nil
}
func Close() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
