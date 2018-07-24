# DB Migrate Database Adaptors

### Postgres Usage

```go
import "bitbucket.org/braindev/dbmigrate/adaptor"

// NewPostgres(db *sql.DB) (*PostgresAdaptor, error)
pgadaptor, err := adaptor.NewPostgres(db)
```

### MySQL Usage

```go
import "bitbucket.org/braindev/dbmigrate/adaptor"

// NewMySQL(db *sql.DB) (*MySQLAdaptor, error)
mysqladaptor, err := adaptor.NewMySQL(db)
```

### Custom Adaptor

To create a custom adaptor implement the following interface:

```go
type VendorAdaptor interface {
	GetAppliedMigrationsOrderedAsc() ([]string, error)
	ApplyMigration(pair MigrationPair) error
	RollbackMigration(pair MigrationPair) error
}
```