# DB Migrate Database Adapters

### Postgres Usage

```go
import "bitbucket.org/braindev/dbmigrate/adapter"

// NewPostgres(db *sql.DB) (*PostgresAdapter, error)
pgadapter, err := adapter.NewPostgres(db)
```

### MySQL Usage

```go
import "bitbucket.org/braindev/dbmigrate/adapter"

// NewMySQL(db *sql.DB) (*MySQLAdapter, error)
mysqladapter, err := adapter.NewMySQL(db)
```

### Custom Adapter

To create a custom adapter implement the following interface:

```go
type VendorAdapter interface {
	GetAppliedMigrationsOrderedAsc() ([]string, error)
	ApplyMigration(pair MigrationPair) error
	RollbackMigration(pair MigrationPair) error
}
```
