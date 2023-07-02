FROM golang:1.20-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o main .

FROM alpine:latest

COPY --from=build /app/main /app/main

WORKDIR /app

EXPOSE 8080

CMD ["./main"]