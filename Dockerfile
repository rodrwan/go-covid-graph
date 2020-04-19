FROM golang:alpine as builder

WORKDIR /app

COPY go.mod .
COPY go.sum .
COPY pkg pkg
COPY main.go .

RUN CGO_ENABLED=0 GOOS=linux go build -a -o bin/server -tags netgo -installsuffix netgo

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /usr/bin

RUN mkdir templates assets

COPY templates templates
COPY --from=builder /app/bin/server .

EXPOSE 5000

CMD ["/bin/sh", "-l", "-c", "server"]

