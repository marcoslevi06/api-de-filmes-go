FROM golang:1.26.5-alpine AS builder

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /app/api ./cmd/api

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/api ./api
COPY --from=builder /src/movies.json ./movies.json

EXPOSE 8080

ENTRYPOINT ["./api"]
