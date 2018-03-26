FROM golang:1.10.0 as builder

WORKDIR /go/src/github.com/asciiu/gomo/api

COPY . .

RUN go get -u github.com/golang/dep/cmd/dep
RUN dep init 
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build


FROM alpine:latest

#RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/asciiu/gomo/api .

EXPOSE 5000

CMD ["./api"]