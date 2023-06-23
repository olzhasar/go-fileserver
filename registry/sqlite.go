package registry

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
)

type SQLiteRegistry struct {
	db *sql.DB
}

func (r *SQLiteRegistry) Record(token, fileName string) error {
	_, err := r.db.Exec("INSERT INTO files (token, filename) VALUES (?, ?)", token, fileName)
	return err
}

func (r *SQLiteRegistry) Get(token string) (fileName string, ok bool) {
	err := r.db.QueryRow("SELECT filename FROM files WHERE token = ?", token).Scan(&fileName)
	if err == sql.ErrNoRows {
		return fileName, false
	}
	return fileName, true
}

func (r *SQLiteRegistry) Has(token string) bool {
	var exists bool
	r.db.QueryRow("SELECT 1 FROM FILES WHERE token = ?", token).Scan(&exists)
	return exists
}

func (r *SQLiteRegistry) Clear() {
	_, err := r.db.Exec("DELETE FROM files;")
	if err != nil {
		panic(err)
	}
}

func (r *SQLiteRegistry) Close() {
	r.db.Close()
}

func NewSQLiteRegistry(dbPath string) (Registry, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return &SQLiteRegistry{}, err
	}

	db.Exec(`
CREATE TABLE IF NOT EXISTS files(
id INTEGER NOT NULL PRIMARY KEY,
token VARCHAR(24),
filename VARCHAR(255)
);
CREATE UNIQUE INDEX IF NOT EXISTS idx_token ON files (token);
`)

	return &SQLiteRegistry{db}, nil
}
