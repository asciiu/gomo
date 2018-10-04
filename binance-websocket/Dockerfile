# We use the official golang image, which contains all the 
# correct build tools and libraries. Notice `as builder`,
# this gives this container a name that we can reference later on. 
FROM golang:1.10.2 as builder

# Set our workdir to our current service in the gopath
WORKDIR /go/src/github.com/asciiu/gomo/binance-websocket

RUN go get -u github.com/golang/dep/cmd/dep

ADD id_rsa /root/.ssh/id_rsa
RUN chmod 700 /root/.ssh/id_rsa
RUN echo "Host github.com\n\tStrictHostKeyChecking no\n" >> /root/.ssh/config
RUN git config --global url.ssh://git@github.com/.insteadOf https://github.com/

ADD binance-websocket /go/src/github.com/asciiu/gomo/binance-websocket
CMD cd /go/src/github.com/asciiu/gomo/binance-websocket && go get github.com/asciiu/gomo && go build -o /gomo

# Create a dep project, and run `ensure`, which will pull in all 
# of the dependencies within this directory.
RUN dep ensure

# Build the binary, with a few flags which will allow
# us to run this binary in Alpine. 
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o binance-websocket 


# Here we're using a second FROM statement, which is strange,
# but this tells Docker to start a new build process with this
# image.
FROM alpine:latest

# Security related package, good to have.
RUN apk --no-cache add ca-certificates

# Same as before, create a directory for our app.
RUN mkdir /app
WORKDIR /app

# Here, instead of copying the binary from our host machine,
# we pull the binary from the container named `builder`, within
# this build context. This reaches into our previous image, finds
# the binary we built, and pulls it into this container. Amazing!
COPY --from=builder /go/src/github.com/asciiu/gomo/binance-websocket/binance-websocket .

# Run the binary as per usual! This time with a binary build in a
# separate container, with all of the correct dependencies and
# run time libraries.
CMD ["./binance-websocket"]