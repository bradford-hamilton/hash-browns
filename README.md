# HashBrowns
1. Give me your password
2. Profit

## Development

### Environment
Be sure to add the following environment variables which you can point locally or get from a team member for stage/prod:
```
HB_PORT
HB_DB_HOST
HB_DB_PORT
HB_DB_USER
HB_DB_NAME
```

Example for local development - place this inside your `~/.bash_profile`:
```
# HashBrowns env variables
export HB_PORT=4000
export HB_DB_HOST=localhost
export HB_DB_PORT=5432
export HB_DB_USER=bradford
export HB_DB_NAME=hash_browns_db
```

### Database
Ensure you have PostgreSQL on your machine, and run the `migration.sql` file to create and migrate the db

### Hot Reloading
Realize https://github.com/oxequa/realize makes development a little nicer by hot reloading our server for us
```
go get github.com/oxequa/realize
```
Then from the root of the project run:
```
realize start
```

Or feel free to run/build `cmd/server/main.go`

### Testing
```
go test ./...
```
Or for a little more clarity with some color:
```
go test -v ./... | sed ''/PASS/s//$(printf "\033[32mPASS\033[0m")/'' | sed ''/FAIL/s//$(printf "\033[31mFAIL\033[0m")/''
```

### GoDoc
Run:
```
godoc -http=:3000
```
and visit `http://localhost:3000/pkg/github.com/bradford-hamilton/hash-browns` for documentation