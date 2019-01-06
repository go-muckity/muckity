FROM golang:1.11.4 as builder
RUN GIT_TERMINAL_PROMPT=1 go get "github.com/tsal/muckity"
WORKDIR /go/src/github.com/go-muckity/muckity/
RUN go get ./...
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /go/src/github.com/go-muckity/muckity/muckity .
CMD ["./muckity"]
