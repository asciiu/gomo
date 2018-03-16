# gomo
A Go port of the beloved fomo api.

### Prerequisites
Install goose. 
$ go get -u github.com/pressly/goose/cmd/goose
https://github.com/pressly/goose

### Getting Started 
Apply the migrations from the "database/goose" directory.
$ goose postgres "user=postgres dbname=gomo_dev sslmode=disable" up

Clean DB
$ goose postgres "user=postgres dbname=gomo_dev sslmode=disable" down-to 0 
