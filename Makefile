BIN_NAME=ciak
BIN_OUTPUT=dist/${BIN_NAME}

fmt:
	go fmt ./...

deps:
	go mod vendor
	go mod verify

build: fmt deps
	go build -o ${BIN_OUTPUT}


DOCKER_IMAGE=garugaru/ciak
DOCKER_IMAGE_ARM=${DOCKER_IMAGE}:armhf
COMPOSE=docker/docker-compose.yml
VERSION=$(shell git rev-parse --short HEAD)
DOCKERFILE_ARMHF=Dockerfile.armhf

docker-up:
	docker-compose -f ${COMPOSE} up

docker-upd:
	docker-compose -f ${COMPOSE} up -d

docker-down:
	docker-compose -f ${COMPOSE} down

docker-build:
	docker build -t ${DOCKER_IMAGE}:latest -t ${DOCKER_IMAGE}:${VERSION} .

docker-push: docker-build
	docker push ${DOCKER_IMAGE}:${VERSION}
	docker push ${DOCKER_IMAGE}:latest

docker-build-arm:
	docker build -t ${DOCKER_IMAGE}:arm-latest -f ${DOCKERFILE_ARMHF} .

docker-push-arm:
	docker build -t ${DOCKER_IMAGE}:arm-latest -f ${DOCKERFILE_ARMHF} .
	docker push ${DOCKER_IMAGE}:arm-latest

docker-push-all: docker-push docker-push-arm
	docker manifest create ${DOCKER_IMAGE}:latest ${DOCKER_IMAGE}:latest ${DOCKER_IMAGE}:arm-latest
	docker manifest annotate ${DOCKER_IMAGE}:latest ${DOCKER_IMAGE}:arm-latest --os linux --arch arm
	docker manifest push ${DOCKER_IMAGE}:latest


