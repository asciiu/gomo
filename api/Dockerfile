FROM golang:1.10.0 as builder

WORKDIR /go/src/github.com/asciiu/gomo/api

COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN go get -u github.com/pressly/goose/cmd/goose
RUN go get -u github.com/satori/go.uuid
RUN dep init 
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/asciiu/gomo/api .

CMD ["./vessel-service"]