# DBMigrate

A basic database migration library.  Status: ALPHA - API may change

The algorithm used works as follows:

To apply one/all migration(s)

- find all possible migrations
- find all previously applied migrations
- apply the next or all the migrations (depending on if one migration is to be applied or all should be applied) in alphanumeric order by migration version that have not been previously applied

To rollback one version

- find all possible migrations
- find all previously applied migrations
- find the highest migration version already applied that is in the list of all possible migrations

### Installation

```go get -u bitbucket.org/braindev/dbmigrate```

### Example Usage

```go
// error handling omitted
db, err := sql.Open("postgres", "postgres://user@localhost/test?sslmode=disable")
adapter, err := adapter.NewPostgres(db)
storage := storage.NewFileStorage("/path/to/migrations/directory")
/*
assume a directory with with the following:
/path/to/migrations/directory/
  0001_migration_name_apply.sql
  0001_migration_name_rollback.sql
  0002_another_migration_name_apply.sql
  0002_another_migration_name_rollback.sql
*/
m, err := dbmigrate.New(adapter, storage)
err = m.ApplyAll()
// or
err = m.ApplyOne()
// or
err = m.RollbackLatest()
```

Both migration storage and database adapters are pluggable.  Their interfaces are simple:

```go
type Storage interface {
	GetMigrationPairs() ([]MigrationPair, error)
}
```

[Included Storage Adapters](./storage/README.md)

```go
type VendorAdapter interface {
	GetAppliedMigrationsOrderedAsc() ([]string, error)
	ApplyMigration(pair MigrationPair) error
	RollbackMigration(pair MigrationPair) error
}
```

[Included Database Adapters](./adapter/README.md)

### Goals

- Simplicity
- Modularity

### Non-Goals

- ORM-like/DSL features

### Some Alternatives

- [gondolier](https://github.com/emvicom/gondolier)
- [goose](https://github.com/steinbacher/goose)
- [gormigrate](https://github.com/go-gormigrate/gormigrate)
- [migrate](https://github.com/golang-migrate/migrate)
- [sql-migrate](https://github.com/rubenv/sql-migrate)

### TODO

- Add more database adapters?
- Add more storage options?
- More tests
- Force migration at a specific version?
- Add a changelog
