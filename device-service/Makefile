build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/asciiu/gomo/device-service proto/device/device.proto
	docker build -t device-service .
run: 
	docker run -p 5002:5000 \
		-e DB_URL="postgres://postgres@docker.for.mac.host.internal:5432/gomo_dev?&sslmode=disable" \
		-e MICRO_SERVER_ADDRESS=:5002 \
		-e MICRO_REGISTRY=mdns device-service