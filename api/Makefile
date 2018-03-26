build:
	docker build -t gomo-api .

run:
	docker run -p 5000:5000 \
		-e MICRO_SERVER_ADDRESS=:50051 \
		-e MICRO_REGISTRY=mdns gomo-api 