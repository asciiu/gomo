build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/asciiu/gomo/order-service proto/order/order.proto
	docker build -t order-service .
run: 
	docker run -p 5004:5000 \
		-e DB_URL="postgres://postgres@docker.for.mac.host.internal:5432/gomo_dev?&sslmode=disable" \
		-e MICRO_SERVER_ADDRESS=:5002 \
		-e MICRO_REGISTRY=mdns order-service