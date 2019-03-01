package adapter

import (
	"database/sql"

	"github.com/braindev/dbmigrate"
)

// PostgresAdapter is a PostgreSQL DBMigrate adapter
type PostgresAdapter struct {
	db queryExecutor
}

// NewPostgres creates a new PostgresAdapter
func NewPostgres(db *sql.DB) (*PostgresAdapter, error) {
	a := &PostgresAdapter{
		db: db,
	}
	return a, a.createMigrationsTable()
}

// GetAppliedMigrationsOrderedAsc returns an ordered slice of string versions
// of migrations that have been previously applied
func (a *PostgresAdapter) GetAppliedMigrationsOrderedAsc() ([]string, error) {
	const query = `SELECT "version" FROM "dbmigrations" ORDER BY "version" ASC`
	rows, err := a.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var version string
	versions := []string{}
	for rows.Next() {
		err = rows.Scan(&version)
		if err != nil {
			return nil, err
		}
		versions = append(versions, version)
	}
	return versions, nil
}

// createMigrationsTable conditionally creates the migrations table if it
// doesn't yet exist
func (a *PostgresAdapter) createMigrationsTable() error {
	const query = `CREATE TABLE IF NOT EXISTS "dbmigrations"("version" varchar NOT NULL PRIMARY KEY)`
	_, err := a.db.Exec(query)
	return err
}

// ApplyMigration applies the specified migration
func (a *PostgresAdapter) ApplyMigration(pair dbmigrate.MigrationPair) error {
	_, err := a.db.Exec(pair.ApplyBody)
	if err != nil {
		return err
	}
	_, err = a.db.Exec(`INSERT INTO "dbmigrations" ("version") VALUES ($1)`, pair.Version)
	return err
}

// RollbackMigration rolls back the specifified migration
func (a *PostgresAdapter) RollbackMigration(pair dbmigrate.MigrationPair) error {
	_, err := a.db.Exec(pair.RollbackBody)
	if err != nil {
		return err
	}
	_, err = a.db.Exec(`DELETE FROM "dbmigrations" WHERE "version" = $1`, pair.Version)
	return err
}
