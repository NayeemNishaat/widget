**To overcome the issue about not being able to create new methods on a struct in a different package we can create a config file in lib package and create a initConfig function that we will call from main function to initialize app config and we will import this init function to get the app config that we might need in handlers.**

# Migration

## Create

- /Users/labyrinth/.go/bin/migrate create -ext sql -dir ./api/migration -seq mg

## Run Migration Up

- /Users/labyrinth/.go/bin/migrate -path ./api/migration -database pgx5://localhost:5432/ecom -verbose up
<!-- /Users/labyrinth/.go/bin/migrate -path ./api/migration -database "postgres://localhost:5432/ecom?sslmode=disable" -verbose up -->

## Run Migration Up

- /Users/labyrinth/.go/bin/migrate -path ./api/migration -database pgx5://localhost:5432/ecom -verbose down

## Resolve Migration Error

- migrate -path database/migration/ -database "pgx5://username:secretkey@localhost:5432/database_name?sslmode=disable" force <VERSION>
<!-- /Users/labyrinth/.go/bin/migrate -path ./api/migration -database "postgres://localhost:5432/ecom?sslmode=disable" force 1 -->

# Makefile

- migration_up: migrate -path database/migration/ -database "pgx5://username:secretkey@localhost:5432/database_name?sslmode=disable" -verbose up

- migration_down: migrate -path database/migration/ -database "pgx5://username:secretkey@localhost:5432/database_name?sslmode=disable" -verbose down

- migration_fix: migrate -path database/migration/ -database "pgx5://user:password@host:port/dbname?query" force VERSION

**Run `make migration_up`, `make migration_down`, `make migration_fix`**

`make migration_fix v=2` # v should be any last successful migration version (current -1 is preffered).

# Install Go Software

- Download the source
- Go to the main file dir and run `go install .` or `go install -tags 'pgx5' .`
- The binary file will be available inside go path.

# Build

Default build location will be in the current dir

## Windows

- GOOS=windows go build
- GOOS=windows go build main.go

## Mac

- GOOS=darwin go build
- GOOS=darwin go build main.go
- go build -tags="mysql sqlite3 postgres mongodb"
