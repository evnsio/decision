FROM golang:1.13.7 AS builder
WORKDIR /go/src/github.com/evnsio/decision/
COPY ./ .
RUN go get ./...
RUN env GOOS=linux GOARCH=amd64 go build ./cmd/decision/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/evnsio/decision/decision .
ENTRYPOINT ["./decision"]
