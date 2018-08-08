# DB Migrate Storage Adaptors

### File Storage

File storage allows migrations to be stored in files.  All migrations must be
in one directory.  Migrations must be in pairs with one being an "apply"
migration and the other being a "rollback" migration.  An "apply" migration
consists of the changes to put the database into a new state.  The "rollback"
migration consists of the changes to get the database back to the state prior
to the "rollack" migration.  ("apply" and "rollback" as opposed to "up" and
"down" are used because how they sort alphabeterically). The migrations must
have versions.  A version is a string is alphanumerically sortable.  In the
case the file storage each migration file name should match the format
described by this regex:

```^([\d\-_]+:)[\-_]([\w\-]+)[\-_](apply|rollback)(\.\w+)?$```

The file name consists of

- the version (numbers, "_", ":", and "-" characters)
- a "-" or "_" character
- a name for the migration (numbers, alpha, "_", and "-" characters)
- a "-" or "_" character
- "apply" or "rollback"
- and optional extension

Examples of possible file names:

```
/path/to/migrations/directory
  0001_migration_name_apply.sql
  0001_migration_name_rollback.sql
```

```
/path/to/migrations/directory
  2020-10-01-15:35:35-migration-name-apply
  2020-10-01-15:35:35-migration-name-rollback
```

Because the versions are strings and sorted alpha-numerically "1" and "10" come
before "2".  This is a **bad** example of a migrations directory.

```
/path/to/migrations/directory
  1-migration-name-apply.sql
  2-migration-name-rollback.sql
  ...
  10-migration-name-apply.sql
  10-migration-name-rollback.sql
```

### Custom Migration Storage

To implement a custom migration storeage simply imlement this interface:

```go
type Storage interface {
	GetMigrationPairs() ([]MigrationPair, error)
}
```

A `MigrationPair` is defined as

```go
type MigrationPair struct {
	Version      string
	Name         string
	ApplyBody    string
	RollbackBody string
}
```