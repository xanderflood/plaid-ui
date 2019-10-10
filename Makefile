.PHONY: build-docker build-local docker local

build-docker:
	GOOS=linux go build -o build/api/api ./cmd/api
	docker build build/api -t xanderflood/plaid-ui-api

build-local:
	go build -o build/api/api ./cmd/api

docker: build-docker
	docker run --publish 8000:8000 --env-file .docker.env xanderflood/plaid-ui-public-api

local: build-local
	cd build/api && godotenv -f ../../.env ./api