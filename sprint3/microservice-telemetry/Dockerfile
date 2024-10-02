FROM golang:1.22

WORKDIR ${GOPATH}/telemetry/
COPY . ${GOPATH}/telemetry/

RUN go build -o /build ./internal/cmd \
    && go clean -cache -modcache

EXPOSE 8081

CMD ["/build"]