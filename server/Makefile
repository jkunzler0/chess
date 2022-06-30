
DOCKER_IMAGE := chess-server

## Build the docker image
docker-build:
	docker build -t $(DOCKER_IMAGE) .

## Run the docker image
docker-run:
	docker run -dp 8080:8080 $(DOCKER_IMAGE)

## Stop the docker image
docker-stop:
	docker stop $(docker ps -a -q)

## Remove the docker image
docker-rm:
	docker rm $(docker ps -a -q)
