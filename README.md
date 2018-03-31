# gomo
A Go port of the beloved fomo api.

### Prerequisites
Install goose. 
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

### Getting Started 
Apply the migrations from the "migrations" directory.
```
$ goose postgres "user=postgres dbname=gomo_dev sslmode=disable" up
```
Clean DB
```
$ goose postgres "user=postgres dbname=gomo_dev sslmode=disable" down-to 0 
```

### Testing 
Apply DB schema to test database. Create dB gomo_test if it does not exist. 

```
$ goose postgres "user=postgres dbname=gomo_test sslmode=disable" up
```


### Generating the API docs
From within the /api project 
```
$ swagger generate spec -o ./swagger.json --scan-models
$ swagger serve -F=swagger swagger.json
```

