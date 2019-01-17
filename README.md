## HashBrowns
1. Give me your password
2. Profit

### Development

#### Environment
Be sure to add the following environment variables which you can point locally or get from a team member for stage/prod:
```
HB_PORT
HB_DB_HOST
HB_DB_PORT
HB_DB_USER
HB_DB_NAME
```

#### Database
Ensure you have PostgreSQL on your machine, and run the `migration.sql` file to create and migrate the db

#### Hot Reloading
Realize https://github.com/oxequa/realize makes development a little nicer by hot reloading our server for us
```
go get github.com/oxequa/realize
```
Then from the root of the project run:
```
realize start
```

Or feel free to run/build `cmd/server/main.go`