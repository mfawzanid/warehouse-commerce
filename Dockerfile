FROM golang:1.23-alpine as Build

ARG TARGET=app

WORKDIR /app

COPY . .

RUN GOPATH= go build -o /app/main cmd/$TARGET/main.go

FROM alpine:latest

ARG TARGET

WORKDIR /root

COPY --from=Build /app/main .

EXPOSE 1323

ENTRYPOINT ["./main"]