FROM golang:1.13-alpine as builder

WORKDIR /app

COPY . .

# Install gcc
RUN apk update && \
    apk add gcc libc-dev

RUN go build

FROM alpine:latest

RUN apk update && \
    apk add libc-dev

WORKDIR /app

COPY --from=builder /app/calorie-tracker .

CMD [ "/app/calorie-tracker" ]