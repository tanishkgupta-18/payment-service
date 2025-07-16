package migrations

import (
	"gofr.dev/pkg/gofr/migration"
)

const createHealthCheckTableSQL = `CREATE TABLE IF NOT EXISTS health_check (
    id INT AUTO_INCREMENT PRIMARY KEY,
    status VARCHAR(10) NOT NULL
);`

func createHealthCheckTable() migration.Migrate {
	return migration.Migrate{
		UP: func(d migration.Datasource) error {
			_, err := d.SQL.Exec(createHealthCheckTableSQL)
			if err != nil {
				return err
			}
			return nil
		},
	}
}
