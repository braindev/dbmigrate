# DBMigrate

A basic database migration library.  Status: ALPHA - API may change

The algorithm used works as follows:

To migrate

### Installation

```go get -u bitbucket.org/braindev/dbmigrate```

### Example Usage

```go
// error handling omitted
db, err := sql.Open("postgres", "postgres://user@localhost/test?sslmode=disable")
adaptor, err := adaptor.NewPostgres(db)
storage := storage.NewFileStorage("/path/to/migrations/directory")
m, err := dbmigrate.New(adaptor, storage)
err = m.ApplyAll()
# or
err = m.ApplyOne()
# or
err = m.RollbackLatest()
```

Both migration storage and database adaptors are pluggable.  Their interfaces are simple:

```go
type Storage interface {
	GetMigrationPairs() ([]MigrationPair, error)
}
```

[Included Storage Adaptors](./storage/README.md)

```go
type VendorAdaptor interface {
	GetAppliedMigrationsOrderedAsc() ([]string, error)
	ApplyMigration(pair MigrationPair) error
	RollbackMigration(pair MigrationPair) error
}
```

[Included Database Adaptors](./adaptor/README.md)

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

- Add more database adaptors?
- Add more storage options?
- More tests
- Force migration at a specific version?
- Add a changelog