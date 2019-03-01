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
- run the rollback of that migration

### Installation

```go get -u github.com/braindev/dbmigrate```

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

### CLI Example

There's no CLI distributed with dbmigrate current but a simple CLI might look like:

```go
package main

import (
	"database/sql"
	"flag"
	"fmt"
	"os"

	"github.com/braindev/dbmigrate"
	"github.com/braindev/dbmigrate/adapter"
	"github.com/braindev/dbmigrate/storage"

	_ "github.com/lib/pq"
)

var (
	applyOne = flag.Bool("apply-one", false, "apply the next migration not yet applied in alphanumeric order")
	applyAll = flag.Bool("apply-all", false, "apply all migrations not yet applied in alphanumeric order by version")
	rollback = flag.Bool("rollback", false, "rollback the last applied migration")
	dir      = flag.String("directory", "migrations", "the path to the migrations directory")
	help     = flag.Bool("help", false, "usage information")
)

func main() {
	flag.Parse()
	verifyFlags()

	if *help {
		flag.Usage()
		os.Exit(0)
	}

	db, err := sql.Open("postgres", os.Getenv("DB_URL")) // example DB_URL=postgres://user@localhost/db
	if err != nil {
		panic(err)
	}
	adapter, err := adapter.NewPostgres(db)
	if err != nil {
		panic(err)
	}
	storage := storage.NewFileStorage(*dir)
	m, err := dbmigrate.New(adapter, storage)
	if err != nil {
		panic(err)
	}

	if *applyAll {
		fmt.Println("Applying migrations")
		err = m.ApplyAll()
	} else if *applyOne {
		fmt.Println("Applying one migration")
		err = m.ApplyOne()
	} else if *rollback {
		err = m.RollbackLatest()
	}

	if err != nil {
		fmt.Println("Error:", err)
		os.Exit(-1)
	}
}

func verifyFlags() {
	actionCount := 0
	if *applyOne {
		actionCount++
	}
	if *applyAll {
		actionCount++
	}
	if *rollback {
		actionCount++
	}
	if actionCount > 1 {
		fmt.Println("Cannot choose more than one action.  Please choose one of apply-one, apply-all, or rollback")
		os.Exit(-1)
	}
}
```


### TODO

- Add more database adapters?
- Add more storage options?
- More tests
- Force migration at a specific version?
- Add a changelog
