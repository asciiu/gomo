build:
	docker build -t gomo-api .

run:
	docker run -p 5000:5000 \
		-e GOMO_JWT=:jyVadHV5Z+tihp4ptuEvD2xXu6iWrFaSebYovtmzl0Cg4RDi1udp7qyme8ZZQl+A \
		-e DB_URL="postgres://postgres@docker.for.mac.host.internal:5432/gomo_dev?&sslmode=disable" \
		-e MICRO_SERVER_ADDRESS=:5001 \
		-e MICRO_REGISTRY=mdns gomo-api