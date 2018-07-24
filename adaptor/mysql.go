package adaptor

import (
	"database/sql"

	"bitbucket.org/braindev/dbmigrate"
)

// MySQLAdaptor is a PostgreSQL DBMigrate adaptor
type MySQLAdaptor struct {
	db *sql.DB
}

// NewMySQL creates a new PostgresAdaptor
func NewMySQL(db *sql.DB) (*MySQLAdaptor, error) {
	a := &MySQLAdaptor{
		db: db,
	}
	return a, a.createMigrationsTable()
}

// GetAppliedMigrationsOrderedAsc returns an ordered slice of string versions
// of migrations that have been previously applied
func (a *MySQLAdaptor) GetAppliedMigrationsOrderedAsc() ([]string, error) {
	const query = "SELECT `version` FROM `dbmigrations` ORDER BY `version` ASC"
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
func (a *MySQLAdaptor) createMigrationsTable() error {
	const query = "CREATE TABLE IF NOT EXISTS `dbmigrations`(`version` varchar NOT NULL PRIMARY KEY)"
	_, err := a.db.Exec(query)
	return err
}

// ApplyMigration applies the specified migration
func (a *MySQLAdaptor) ApplyMigration(pair dbmigrate.MigrationPair) error {
	_, err := a.db.Exec(pair.ApplyBody)
	if err != nil {
		return err
	}
	_, err = a.db.Exec("INSERT INTO `dbmigrations` (`version`) VALUES (?)", pair.Version)
	return err
}

// RollbackMigration rolls back the specifified migration
func (a *MySQLAdaptor) RollbackMigration(pair dbmigrate.MigrationPair) error {
	_, err := a.db.Exec(pair.RollbackBody)
	if err != nil {
		return err
	}
	_, err = a.db.Exec("DELETE FROM `dbmigrations` WHERE `version` = ?", pair.Version)
	return nil
}
