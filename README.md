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

### Migrate the Postgres DB schema 
Apply the migrations from the "migrations" directory.
```
$ goose postgres "user=postgres dbname=gomo_dev sslmode=disable" up
```

Note: When running from docker-compose up you need to migrate the dockerized postgres DB via:
```
$ goose postgres "user=fomo dbname=fomo_dev sslmode=disable port=6432 password=fomornd" up
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
$ swagger generate spec -o ./fomo-swagger.json --scan-models
$ swagger serve -F=swagger fomo-swagger.json
```

### Deploying
From localhost using docker-machine you first need to create the ec2 instances:

```
docker-machine create --driver amazonec2 --amazonec2-region us-west-1 fomo-stage
```

Set the docker env:
```
eval $(docker-machine env fomo-stage)
```

Deploy via compose build and up. 
```
docker-compose build
docker-compose -f docker-compose.yml -f docker-compose.stage.yml up -d
```

