FROM golang:1.20.3 AS builder
COPY . $GOPATH/src/app
WORKDIR $GOPATH/src/app
RUN go get -d -v
RUN GO111MODULE=on CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CGO_CFLAGS="-g -O2 -Wno-return-local-addr" go build -o $GOPATH/src/app/app.bin

FROM alpine
ENV TZ Europe/Moscow
COPY --from=builder /go/src/app/app.bin /app.bin
EXPOSE 8443/tcp
EXPOSE 1812/tcp
VOLUME ["/certs","/config"]
CMD ["/app.bin"]