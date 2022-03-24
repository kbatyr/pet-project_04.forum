run:
	@echo running server...
	go run ./cmd/web
	@echo

build:
	@echo building project...
		go build -o forum
	@echo

docker:
	@echo building image...
	docker build --tag docker-forum .
	@echo

	@echo image list
	docker images
	@echo

	@echo creating and running container...
	docker run -d -p 8080:8080 --name forum-container docker-forum
	@echo
	
	@echo container list
	docker ps -a
	@echo

prune:
	@echo delete unused images and containers...
	docker image prune --filter "until=24h"
	docker container prune --filter "until=24h"
	@echo

stop:
	@echo stoping container...
	docker stop forum-container
	@echo

	@echo removing container...
	docker rm forum-container
	@echo

	@echo removing image...
	docker image rm docker-forum
	docker image rm 80d9a75ccb38
	@echo

.DEFAULT_GOAL:= run