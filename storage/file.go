package storage

import (
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"bitbucket.org/braindev/dbmigrate"
)

var (
	migratationFileNameRe = regexp.MustCompile(`^([\d\-_]+)[\-_]([\w\-]+)[\-_](apply|rollback)(\.\w+)?$`)
)

// FileStorage ...
type FileStorage struct {
	directoryPath string
}

// NewFileStorage creates & return a new FileStore struct
func NewFileStorage(path string) *FileStorage {
	return &FileStorage{
		directoryPath: path,
	}
}

// GetMigrationPairs ...
func (f *FileStorage) GetMigrationPairs() ([]dbmigrate.MigrationPair, error) {
	if info, err := os.Stat(f.directoryPath); err != nil {
		return nil, err
	} else if !info.IsDir() {
		return nil, errors.New(f.directoryPath + " is not a directory")
	}
	files, err := filepath.Glob(filepath.Join(f.directoryPath, "*"))
	if err != nil {
		return nil, err
	}
	migrations := map[string]dbmigrate.MigrationPair{}
	for _, file := range files {
		matches := migratationFileNameRe.FindStringSubmatch(filepath.Base(file))
		if len(matches) >= 4 {
			version := matches[1]
			var pair dbmigrate.MigrationPair
			if _, ok := migrations[version]; ok {
				pair = migrations[version]
			} else {
				pair = dbmigrate.MigrationPair{
					Version: version,
					Name:    matches[2],
				}
			}
			if matches[3] == "apply" {
				if b, err := ioutil.ReadFile(file); err == nil {
					pair.ApplyBody = string(b)
				} else {
					return nil, err
				}
			} else {
				if b, err := ioutil.ReadFile(file); err == nil {
					pair.RollbackBody = string(b)
				} else {
					return nil, err
				}
			}
			migrations[version] = pair
		}
	}
	migrationFiles := []dbmigrate.MigrationPair{}
	for v, m := range migrations {
		if migrations[v].ApplyBody == "" {
			return nil, errors.New(`apply migration not found for version ` + v)
		}
		if migrations[v].RollbackBody == "" {
			return nil, errors.New(`rollback migration not found for version ` + v)
		}
		migrationFiles = append(migrationFiles, m)
	}

	return migrationFiles, err
}
