# Stage 1: Build
FROM golang:1.26-alpine3.23 AS builder

WORKDIR /app

COPY . .
RUN go mod download && go build -o sourcescore cmd/app/main.go

# Stage 2: Run
FROM alpine:3.23

WORKDIR /app

COPY --from=builder /app/sourcescore ./sourcescore

EXPOSE 8080

CMD ["./sourcescore"]