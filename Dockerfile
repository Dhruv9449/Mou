# Builder
FROM golang:1.19.3-alpine3.15 AS builder

RUN apk update \
    && apk --no-cache --update add build-base git

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o bin/mou


# Runner
FROM alpine:3.15 as runner

WORKDIR /app

COPY --from=build /app/bin/mou ./

RUN chmod +x mou

CMD ["./mou"]