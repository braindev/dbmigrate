package dbmigrate

import (
	"fmt"
	"sort"
)

// VendorAdapter is the set of functionality needed to implement a compatibility with a database
// vendor
type VendorAdapter interface {
	GetAppliedMigrationsOrderedAsc() ([]string, error)
	ApplyMigration(pair MigrationPair) error
	RollbackMigration(pair MigrationPair) error
}

// Storage is an interface for where the migrations are stored
type Storage interface {
	GetMigrationPairs() ([]MigrationPair, error)
}

// MigrationPair stores a pair of apply/rollback migrations
type MigrationPair struct {
	Version      string
	Name         string
	ApplyBody    string
	RollbackBody string
}

// DBMigrate stores all the configuration needed to run migrations
type DBMigrate struct {
	sortedMigrationPairs []MigrationPair
	migrationPairs       map[string]MigrationPair
	allVersions          []string
	adapter              VendorAdapter
	storage              Storage
	verbose              bool // TODO
}

type migrationPairs []MigrationPair

func (a migrationPairs) Len() int           { return len(a) }
func (a migrationPairs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a migrationPairs) Less(i, j int) bool { return a[i].Version < a[j].Version }

// New returns an initialized DBMigrate instance.
func New(adapter VendorAdapter, storage Storage) (*DBMigrate, error) {
	dbMigrate := DBMigrate{
		adapter: adapter,
		storage: storage,
	}
	if migrationPairs, err := storage.GetMigrationPairs(); err == nil {
		dbMigrate.sortedMigrationPairs = migrationPairs
	} else {
		return nil, err
	}
	sort.Sort(migrationPairs(dbMigrate.sortedMigrationPairs))

	dbMigrate.migrationPairs = map[string]MigrationPair{}
	for _, m := range dbMigrate.sortedMigrationPairs {
		dbMigrate.migrationPairs[m.Version] = m
	}

	return &dbMigrate, nil
}

// ApplyAll applies all migrations versions not yet applied in ascending
// alphanumeric order
func (dbMigrate *DBMigrate) ApplyAll() error {
	return dbMigrate.apply(true)
}

// ApplyOne applies the next not yet applied migrations going in ascending
// alphenermic order
func (dbMigrate *DBMigrate) ApplyOne() error {
	return dbMigrate.apply(false)
}

func (dbMigrate *DBMigrate) apply(all bool) error {
	appliedVersions, err := dbMigrate.adapter.GetAppliedMigrationsOrderedAsc()
	if err != nil {
		return err
	}

	// TODO this could be more efficient
	versionsToApply := []string{}
	for version := range dbMigrate.migrationPairs {
		if !stringInSlice(version, appliedVersions) {
			versionsToApply = append(versionsToApply, version)
		}
	}

	sort.Sort(sort.StringSlice(versionsToApply))

	for _, version := range versionsToApply {
		fmt.Println("Applying migration", dbMigrate.migrationPairs[version].Name)
		err = dbMigrate.adapter.ApplyMigration(dbMigrate.migrationPairs[version])
		if err != nil {
			return err
		}
		if !all {
			return nil
		}
	}

	return nil
}

// RollbackLatest rolls back the latest applied migration
func (dbMigrate *DBMigrate) RollbackLatest() error {
	appliedVersions, err := dbMigrate.adapter.GetAppliedMigrationsOrderedAsc()
	if err != nil {
		return err
	}
	for i := len(dbMigrate.sortedMigrationPairs) - 1; i >= 0; i-- {
		if stringInSlice(dbMigrate.sortedMigrationPairs[i].Version, appliedVersions) {
			fmt.Println("rolling back migration", dbMigrate.sortedMigrationPairs[i].Name)
			err = dbMigrate.adapter.RollbackMigration(dbMigrate.sortedMigrationPairs[i])
			if err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func stringInSlice(needle string, haystack []string) bool {
	for _, s := range haystack {
		if s == needle {
			return true
		}
	}
	return false
}
