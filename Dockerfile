FROM golang:1.23-alpine as builder
ARG VERSION_TAG
ARG COMMIT_ID
ARG BUILD_DATE
RUN apk add git make build-base
COPY src /srv/fsmq
WORKDIR /srv/fsmq
ENV GOBIN=/usr/local/bin
ENV CGO_ENABLED=1
RUN go get && go build -ldflags="-X main.Version=$VERSION_TAG -X main.CommitID=$COMMIT_ID -X main.BuildDate=$BUILD_DATE" -buildvcs=false -o fsmq

FROM alpine:latest
COPY --from=builder /srv/fsmq/fsmq /srv/fsmq/fsmq
CMD /srv/fsmq/fsmq
