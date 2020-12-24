FROM golang:alpine as builder
RUN apk add git
COPY src/ /srv/fsmq
WORKDIR /srv/fsmq
ENV GOBIN=/usr/local/bin
RUN go get && go build -o fsmq

FROM alpine:latest
COPY --from=builder /srv/fsmq/fsmq /srv/fsmq/fsmq
CMD /srv/fsmq/fsmq
