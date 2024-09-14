# BASE GO IMAGE
FROM golang:1.22-alpine as builder

RUN mkdir /app

COPY . /app

WORKDIR /app

RUN make build

# BUILD A LIGHT IMAGE
FROM alpine:latest

RUN make run

COPY --from=builder /app/bin/parser /app

CMD [ "/app/parser" ]