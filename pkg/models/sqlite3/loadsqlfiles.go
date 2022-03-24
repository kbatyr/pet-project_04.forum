package sqlite3

import (
	"database/sql"
	"io/ioutil"
	"path/filepath"
	"strings"
)

func LoadSqlFiles(dir string, db *sql.DB) error {

	paths, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return err
	}

	for _, path := range paths {

		file, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		tx, err := db.Begin()
		if err != nil {
			return err
		}

		defer func() {
			tx.Rollback()
		}()

		for _, q := range strings.Split(string(file), ";") {

			q := strings.TrimSpace(q)
			if q == "" {
				continue
			}

			if _, err := tx.Exec(q); err != nil {
				return err
			}
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}
