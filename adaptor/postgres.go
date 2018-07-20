package adaptor

import (
	"database/sql"
	"io/ioutil"

	"bitbucket.org/dbmigrate"
)

// PostgresAdaptor is a PostgreSQL DBMigrate adaptor
type PostgresAdaptor struct {
	db *sql.DB
}

// NewPostgres creates a new PostgresAdaptor
func NewPostgres(db *sql.DB) *PostgresAdaptor {
	return &PostgresAdaptor{
		db: db,
	}
}

// GetAppliedMigrationsOrderedAsc returns an ordered slice of string versions
// of migrations that have been previously applied
func (a *PostgresAdaptor) GetAppliedMigrationsOrderedAsc() ([]string, error) {
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

// CreateMigrationsTable conditionally creates the migrations table if it
// doesn't yet exist
func (a *PostgresAdaptor) CreateMigrationsTable() error {
	const query = `CREATE TABLE IF NOT EXISTS "dbmigrations"("version" varchar NOT NULL PRIMARY KEY)`
	_, err := a.db.Exec(query)
	return err
}

// ApplyMigration applies the specified migration
func (a *PostgresAdaptor) ApplyMigration(pair dbmigrate.MigrationFilePair) error {
	b, err := ioutil.ReadFile(pair.ApplyPath)
	if err != nil {
		return err
	}
	_, err = a.db.Exec(string(b))
	if err != nil {
		return err
	}
	_, err = a.db.Exec(`INSERT INTO "dbmigrations" ("version") VALUES ($1)`, pair.Version)
	return err
}

// RollbackMigration rolls back the specifified migration
func (a *PostgresAdaptor) RollbackMigration(pair dbmigrate.MigrationFilePair) error {
	b, err := ioutil.ReadFile(pair.RollbackPath)
	if err != nil {
		return err
	}
	_, err = a.db.Exec(string(b))
	if err != nil {
		return err
	}
	_, err = a.db.Exec(`DELETE FROM "dbmigrations" WHERE "version" = $1`, pair.Version)
	return nil
}
