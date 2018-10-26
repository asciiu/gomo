# gomo
A Go port of the beloved fomo api.

### Prerequisites
Install goose. This is the database migratin tool. All migrations can be found under
/migrations as incremental SQL scripts.
```
$ go get -u github.com/pressly/goose/cmd/goose
```
https://github.com/pressly/goose

Install go-swagger
```
$ brew tap go-swagger/go-swagger
$ brew install go-swagger
```
Refer to swagger markeup guide here: https://goswagger.io/generate/spec.html

### Migrate the Postgres DB schema 
Apply the migrations from the "migrations" directory.
```
$ goose postgres "user=postgres dbname=gomo_test sslmode=disable" up
```

Clean DB if needed and reapply goose up command above.
```
$ goose postgres "user=postgres dbname=gomo_test sslmode=disable" down-to 0 
```

### Building 
Refer to the Makefile targets: build, stage, etc.

### Deploying 
Refer to the deploment configs. 

### Generating the API docs
From within the /api project 
```
$ swagger generate spec -o ./fomo-swagger.json --scan-models
$ swagger serve -F=swagger fomo-swagger.json
```
