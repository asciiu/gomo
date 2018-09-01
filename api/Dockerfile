FROM golang:1.10.2 as builder

WORKDIR /go/src/github.com/asciiu/gomo/api

# add key so we can pull from private repo 
ADD id_rsa /root/.ssh/id_rsa
RUN chmod 700 /root/.ssh/id_rsa
RUN echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/

ADD api /go/src/github.com/asciiu/gomo/api
CMD cd /go/src/github.com/asciiu/gomo/api && go get github.com/asciiu/gomo && go build -o /gomo


RUN go get -u github.com/golang/dep/cmd/dep
RUN dep ensure -vendor-only
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build


FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/github.com/asciiu/gomo/api .

EXPOSE 5000

CMD ["./api"]