# bug-tracker
Bug Tracking App

## Structure

Not bothering with any framework, will just use the standard library. Go's template engine will be used to handle layouts and pages built with HTMX.

Database is Postgres and pgx is driver used to connect to it. This also uses goose to handle database migrations.

## Links

### Generic Go and HTMX

- https://go.dev/doc/
- https://htmx.org/
- https://gowebexamples.com/

### Package Documentation

- https://pkg.go.dev/html/template
- https://pkg.go.dev/net/http
- https://pkg.go.dev/crypto/sha256
- https://gorilla.github.io/
- https://github.com/jackc/pgx

### Database

- https://www.postgresql.org/
- https://github.com/pressly/goose

Quick note on connecting to local instance

`psql -h localhost -p 5432 -U ${POSTGRES_USER} -d ${POSTGRES_DB}`
