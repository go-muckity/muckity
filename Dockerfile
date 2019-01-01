FROM golang:1.11.4 as builder
WORKDIR /go/src/github.com/tsal/muckity/docker/
COPY game.go .
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/tsal/muckity/docker/app .
CMD ["./app"]
