DRONE_TAG=$(shell git describe --tags --abbrev=0)
DRONE_COMMIT_ID=$(shell git log --format="%H" -n 1)
DRONE_BUILD_DATE=$(shell date "+%Y-%m-%d")

all: fsmq

fsmq: Makefile src/main.go
	podman build \
		--build-arg VERSION_TAG=${DRONE_TAG} \
		--build-arg COMMIT_ID=${DRONE_COMMIT_ID} \
		--build-arg BUILD_DATE=${DRONE_BUILD_DATE} \
		-t docker.io/reactorcoremeltdown/fsmq:${DRONE_TAG} .

login:
	podman login docker.io -u ${DOCKERHUB_USERNAME} -p ${DOCKERHUB_PASSWORD}

release: fsmq login
	pwd
	podman push docker.io/reactorcoremeltdown/fsmq:${DRONE_TAG}
