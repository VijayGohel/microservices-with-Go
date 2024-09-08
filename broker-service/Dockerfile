FROM golang:1.18-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN CGO_ENABLED=0 go build -o broker-app ./cmd/api

RUN chmod +x /app/broker-app

#build tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY --from=builder /app/broker-app /app

CMD [ "/app/broker-app" ]
