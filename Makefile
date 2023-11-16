all: fsmq

fsmq: Makefile src/main.go
	podman build -t docker.io/reactorcoremeltdown/fsmq:latest .

release:
	podman build -t docker.io/reactorcoremeltdown/fsmq:$(TAG) .
	podman push docker.io/reactorcoremeltdown/fsmq:$(TAG)
