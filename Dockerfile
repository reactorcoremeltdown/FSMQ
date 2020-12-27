FROM golang:alpine as builder
RUN apk add git make
COPY . /srv/fsmq
WORKDIR /srv/fsmq
RUN make

FROM alpine:latest
COPY --from=builder /srv/fsmq/src/fsmq /srv/fsmq/fsmq
CMD /srv/fsmq/fsmq
