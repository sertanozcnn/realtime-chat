APP_PATH=./cmd/api
SWAGGER_MAIN=cmd/api/main.go
MONGO_CONTAINER=mongodb
MONGO_VOLUME=mongodb_data

.PHONY: run swagger test fmt tidy mongo-up mongo-down mongo-start mongo-logs mongo-ui-up mongo-ui-down

run:
	go run $(APP_PATH)

swagger:
	go run github.com/swaggo/swag/cmd/swag init -g $(SWAGGER_MAIN) -o docs --parseInternal

test:
	go test ./...

fmt:
	go fmt ./...

tidy:
	go mod tidy

mongo-up:
	docker run -d --name $(MONGO_CONTAINER) -p 27017:27017 -v $(MONGO_VOLUME):/data/db mongo:latest

mongo-down:
	docker stop $(MONGO_CONTAINER)

mongo-start:
	docker start $(MONGO_CONTAINER)

mongo-logs:
	docker logs -f $(MONGO_CONTAINER)

mongo-ui-up:
	docker network inspect realtime-net >/dev/null 2>&1 || docker network create realtime-net
	docker network connect realtime-net $(MONGO_CONTAINER) >/dev/null 2>&1 || true
	docker run -d --name mongo-express --network realtime-net -p 8081:8081 -e ME_CONFIG_MONGODB_URL=mongodb://$(MONGO_CONTAINER):27017 -e ME_CONFIG_BASICAUTH_USERNAME=admin -e ME_CONFIG_BASICAUTH_PASSWORD=admin mongo-express:latest

mongo-ui-down:
	docker stop mongo-express
