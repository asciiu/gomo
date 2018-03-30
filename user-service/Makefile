build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/asciiu/gomo/user-service proto/user/user.proto
	docker build -t user-service .
run: 
	docker run -p 5001:5001 \
		-e DB_URL="postgres://postgres@docker.for.mac.host.internal:5432/gomo_dev?&sslmode=disable" \
		-e MICRO_SERVER_ADDRESS=:5001 \
		-e MICRO_REGISTRY=mdns user-service