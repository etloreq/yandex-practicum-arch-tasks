FROM golang:1.22

WORKDIR ${GOPATH}/devices/
COPY . ${GOPATH}/devices/

RUN go build -o /build ./internal/cmd \
    && go clean -cache -modcache

EXPOSE 8080

CMD ["/build"]