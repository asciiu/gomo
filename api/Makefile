build:
	docker build -t gomo-api .

run:
	docker run -p 5000:5000 \
		-e GOMO_JWT=:jyVadHV5Z+tihp4ptuEvD2xXu6iWrFaSebYovtmzl0Cg4RDi1udp7qyme8ZZQl+A \
		-e DB_PORT=5432 \
		-e DB_HOST=docker.for.mac.host.internal \
		-e DB_NAME=gomo_dev
		gomo-api