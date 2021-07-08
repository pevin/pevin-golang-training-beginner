FROM golang:alpine as builder

RUN apk update

RUN mkdir /app

WORKDIR /app

COPY . .

RUN GOOS=linux GARCH=amd64 go build -ldflags="-w -s" -o /app/engine

FROM alpine:latest as main

RUN mkdir /app

WORKDIR /app

COPY --from=builder /app/engine /app/engine

COPY script/run_cron_entrypoint.sh run_cron_entrypoint.sh

RUN echo '*  *  *  *  *    /app/engine cron > /dev/stdout' > /etc/crontabs/root

EXPOSE 8080
CMD ["/app/engine", "rest"]