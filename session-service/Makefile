build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/asciiu/gomo/session-service proto/session/session.proto
	docker build -t session-service .
run: 
	docker run -p 5001:5001 \
		-e DB_URL="postgres://postgres@docker.for.mac.host.internal:5432/gomo_dev?&sslmode=disable" \
		-e MICRO_SERVER_ADDRESS=:5001 \
		-e MICRO_REGISTRY=mdns session-service