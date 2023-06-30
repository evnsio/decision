FROM golang:1.19 AS builder
WORKDIR /go/src/github.com/evnsio/decision/
COPY ./ .
ENV GOOS=linux
ENV GOARCH=amd64
ENV CGO_ENABLED=0
RUN go get ./...
RUN go build ./cmd/decision/

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/evnsio/decision/decision .
ENTRYPOINT ["./decision"]
