package dbmigrate

import (
	"fmt"
	"testing"
)

func TestNew(t *testing.T) {
	storage := &testStorage{
		migrationPairs: []MigrationPair{
			{
				Name:         "A",
				Version:      "02",
				ApplyBody:    "",
				RollbackBody: "",
			},
			{
				Name:         "B",
				Version:      "01",
				ApplyBody:    "",
				RollbackBody: "",
			},
		},
	}
	adaptor := newTestAdaptor()
	dbm, err := New(adaptor, storage)
	if err != nil {
		t.Error("unexepcted error from New function", err)
	}
	if l := len(dbm.sortedMigrationPairs); l != len(storage.migrationPairs) {
		t.Errorf("expected %d migration pairs got %d", len(storage.migrationPairs), l)
	}
	if dbm.sortedMigrationPairs[0].Version != "01" {
		t.Error("migration pairs should be sorted")
	}
}

func TestApplyAll(t *testing.T) {
	testTable := []struct {
		previousVersions []string
		expectedApplied  []string
	}{
		{
			previousVersions: []string{"002", "003"},
			expectedApplied:  []string{"001", "004", "005"},
		},
		{
			previousVersions: []string{},
			expectedApplied:  []string{"001", "002", "003", "004", "005"},
		},
		{
			previousVersions: []string{"001", "002", "003", "004", "005"},
			expectedApplied:  []string{},
		},
		{
			previousVersions: []string{"002", "003", "004", "005"},
			expectedApplied:  []string{"001"},
		},
	}

	for _, testCase := range testTable {
		testCase := testCase // capture var
		name := fmt.Sprintf("ApplyAll with %v as previous", testCase.previousVersions)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			storage := storageWithMigrations()
			adaptor := newTestAdaptor()
			dbm, _ := New(adaptor, storage)
			adaptor.appliedMigrations = testCase.previousVersions
			dbm.ApplyAll()
			if len(adaptor.appliedMigrationPairs) != len(testCase.expectedApplied) {
				t.Errorf(
					"applied %d migrations but expected %d to be applied",
					len(adaptor.appliedMigrationPairs),
					len(storage.migrationPairs),
				)
			}
			for idx, v := range testCase.expectedApplied {
				if v != adaptor.appliedMigrationPairs[idx].Version {
					t.Errorf(
						"expected version %s applied at index %d but got %s",
						v,
						idx,
						adaptor.appliedMigrationPairs[idx].Version,
					)
				}
			}
		})
	}
}

func TestApplyOne(t *testing.T) {
	testTable := []struct {
		previousVersions []string
		expectedApplied  []string
	}{
		{
			previousVersions: []string{"002", "003"},
			expectedApplied:  []string{"001"},
		},
		{
			previousVersions: []string{},
			expectedApplied:  []string{"001"},
		},
		{
			previousVersions: []string{"001", "002", "003", "004", "005"},
			expectedApplied:  []string{},
		},
		{
			previousVersions: []string{"002", "003", "004", "005"},
			expectedApplied:  []string{"001"},
		},
	}

	for _, testCase := range testTable {
		testCase := testCase // capture var
		name := fmt.Sprintf("ApplyOne with %v as previous", testCase.previousVersions)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			storage := storageWithMigrations()
			adaptor := newTestAdaptor()
			dbm, _ := New(adaptor, storage)
			adaptor.appliedMigrations = testCase.previousVersions
			dbm.ApplyOne()
			if len(adaptor.appliedMigrationPairs) != len(testCase.expectedApplied) {
				t.Errorf(
					"applied %d migrations but expected %d to be applied",
					len(adaptor.appliedMigrationPairs),
					len(storage.migrationPairs),
				)
			}
			for idx, v := range testCase.expectedApplied {
				if v != adaptor.appliedMigrationPairs[idx].Version {
					t.Errorf(
						"expected version %s applied at index %d but got %s",
						v,
						idx,
						adaptor.appliedMigrationPairs[idx].Version,
					)
				}
			}
		})
	}
}

func TestRollbackLatest(t *testing.T) {
	testTable := []struct {
		previousVersions   []string
		expectedRolledback []string
	}{
		{
			previousVersions:   []string{},
			expectedRolledback: []string{},
		},
		{
			previousVersions:   []string{"002"},
			expectedRolledback: []string{"002"},
		},
		{
			previousVersions:   []string{"001", "002", "004"},
			expectedRolledback: []string{"004"},
		},
		{
			previousVersions:   []string{"001", "002", "003", "004", "005"},
			expectedRolledback: []string{"005"},
		},
	}

	for _, testCase := range testTable {
		testCase := testCase // capture var
		name := fmt.Sprintf("ApplyRollbackLatest with %v as previous", testCase.previousVersions)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			storage := storageWithMigrations()
			adaptor := newTestAdaptor()
			dbm, _ := New(adaptor, storage)
			adaptor.appliedMigrations = testCase.previousVersions
			dbm.RollbackLatest()
			if len(adaptor.rolledbackdMigrationPairs) != len(testCase.expectedRolledback) {
				t.Errorf(
					"rolled back %d migrations but expected %d to be rolled back",
					len(adaptor.appliedMigrationPairs),
					len(storage.migrationPairs),
				)
			}
			for idx, v := range testCase.expectedRolledback {
				if v != adaptor.rolledbackdMigrationPairs[idx].Version {
					t.Errorf(
						"expected version %s rolledback at index %d but got %s",
						v,
						idx,
						adaptor.rolledbackdMigrationPairs[idx].Version,
					)
				}
			}
		})
	}

}

type testStorage struct {
	migrationPairs []MigrationPair
}

func (s *testStorage) GetMigrationPairs() ([]MigrationPair, error) {
	return s.migrationPairs, nil
}

func storageWithMigrations() *testStorage {
	storage := &testStorage{
		migrationPairs: []MigrationPair{
			{
				Name:         "004",
				Version:      "004",
				ApplyBody:    "",
				RollbackBody: "",
			},
			{
				Name:         "003",
				Version:      "003",
				ApplyBody:    "",
				RollbackBody: "",
			},
			{
				Name:         "002",
				Version:      "002",
				ApplyBody:    "",
				RollbackBody: "",
			},
			{
				Name:         "005",
				Version:      "005",
				ApplyBody:    "",
				RollbackBody: "",
			},
			{
				Name:         "001",
				Version:      "001",
				ApplyBody:    "",
				RollbackBody: "",
			},
		},
	}
	return storage
}

func newTestAdaptor() *testAdaptor {
	return &testAdaptor{
		appliedMigrations:         []string{}, // previously applied migrations
		appliedMigrationPairs:     []MigrationPair{},
		rolledbackdMigrationPairs: []MigrationPair{},
	}
}

type testAdaptor struct {
	appliedMigrations         []string
	appliedMigrationPairs     []MigrationPair
	rolledbackdMigrationPairs []MigrationPair
}

func (a *testAdaptor) GetAppliedMigrationsOrderedAsc() ([]string, error) {
	return a.appliedMigrations, nil
}

func (a *testAdaptor) ApplyMigration(pair MigrationPair) error {
	a.appliedMigrationPairs = append(a.appliedMigrationPairs, pair)
	return nil
}

func (a *testAdaptor) RollbackMigration(pair MigrationPair) error {
	a.rolledbackdMigrationPairs = append(a.rolledbackdMigrationPairs, pair)
	return nil
}
