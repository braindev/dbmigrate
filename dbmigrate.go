package dbmigrate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

var (
	migratationFileNameRe = regexp.MustCompile(`^(\d+)[\-_]([\w\-]+)[\-_](apply|rollback)\.sql$`)
)

// VendorAdaptor is the set of functionality needed to implement a compatibility with a database
// vendor
type VendorAdaptor interface {
	CreateMigrationsTable() error
	GetAppliedMigrationsOrderedAsc() ([]string, error)
	ApplyMigration(pair MigrationFilePair) error
	RollbackMigration(pair MigrationFilePair) error
}

// MigrationFilePair stores the a pair of apply/rollback migrations
type MigrationFilePair struct {
	ApplyPath    string
	RollbackPath string
	Version      string
	Name         string
}

type migrationFilePairs []MigrationFilePair

func (a migrationFilePairs) Len() int           { return len(a) }
func (a migrationFilePairs) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a migrationFilePairs) Less(i, j int) bool { return a[i].Version < a[j].Version }

// DBMigrate stores all the configuration needed to run migrations
type DBMigrate struct {
	migrationsDirectory  string
	sortedMigrationFiles []MigrationFilePair
	migrationFiles       map[string]MigrationFilePair
	allVersions          []string
	adaptor              VendorAdaptor
	verbose              bool // TODO
}

// New returns an initialized DBMigrate instance.  The migrations directory
// will be scanned and verified.
func New(adaptor VendorAdaptor, dir string) (*DBMigrate, error) {
	dbMigrate := DBMigrate{
		adaptor:             adaptor,
		migrationsDirectory: dir,
	}
	if err := dbMigrate.initialize(); err != nil {
		return nil, err
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
	appliedVersions, err := dbMigrate.adaptor.GetAppliedMigrationsOrderedAsc()
	if err != nil {
		return err
	}

	// TODO this could be more efficient
	versionsToApply := []string{}
	for version := range dbMigrate.migrationFiles {
		if !stringInSlice(version, appliedVersions) {
			versionsToApply = append(versionsToApply, version)
		}
	}

	for _, version := range versionsToApply {
		fmt.Println("Applying migration", dbMigrate.migrationFiles[version].ApplyPath)
		err = dbMigrate.adaptor.ApplyMigration(dbMigrate.migrationFiles[version])
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
	appliedVersions, err := dbMigrate.adaptor.GetAppliedMigrationsOrderedAsc()
	if err != nil {
		return err
	}
	for i := len(dbMigrate.sortedMigrationFiles) - 1; i >= 0; i-- {
		if stringInSlice(dbMigrate.sortedMigrationFiles[i].Version, appliedVersions) {
			fmt.Println("rolling back migration", dbMigrate.sortedMigrationFiles[i].RollbackPath)
			err = dbMigrate.adaptor.RollbackMigration(dbMigrate.sortedMigrationFiles[i])
			if err != nil {
				return err
			}
			return nil
		}
	}

	return nil
}

func (dbMigrate *DBMigrate) initialize() error {
	dir := dbMigrate.migrationsDirectory
	if info, err := os.Stat(dir); err != nil {
		return err
	} else if !info.IsDir() {
		return errors.New(dir + " is not a directory")
	}
	files, err := filepath.Glob(filepath.Join(dir, "*.sql"))
	if err != nil {
		return err
	}
	migrations := map[string]MigrationFilePair{}
	for _, file := range files {
		matches := migratationFileNameRe.FindStringSubmatch(filepath.Base(file))
		if len(matches) == 4 {
			version := matches[1]
			var pair MigrationFilePair
			if _, ok := migrations[version]; ok {
				pair = migrations[version]
			} else {
				pair = MigrationFilePair{
					Version: version,
					Name:    matches[2],
				}
			}
			if matches[3] == "apply" {
				pair.ApplyPath = file
			} else {
				pair.RollbackPath = file
			}
			migrations[version] = pair
		}
	}

	migrationFiles := []MigrationFilePair{}
	for v, m := range migrations {
		if migrations[v].ApplyPath == "" {
			return errors.New(`apply migration not found for version ` + v)
		}
		if migrations[v].RollbackPath == "" {
			return errors.New(`rollback migration not found for version ` + v)
		}
		migrationFiles = append(migrationFiles, m)
	}

	sort.Sort(migrationFilePairs(migrationFiles))

	dbMigrate.sortedMigrationFiles = migrationFiles
	dbMigrate.migrationFiles = migrations

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
