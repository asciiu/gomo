build:
	protoc -I. --go_out=plugins=micro:$(GOPATH)/src/github.com/asciiu/gomo/apikey-service proto/apikey/apikey.proto
	docker build -t apikey-service .
run: 
	docker run -p 5003:5000 \
		-e DB_URL="postgres://postgres@docker.for.mac.host.internal:5432/gomo_dev?&sslmode=disable" \
		-e MICRO_SERVER_ADDRESS=:5002 \
		-e MICRO_REGISTRY=mdns apikey-service