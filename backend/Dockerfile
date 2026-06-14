# syntax=docker/dockerfile:1

FROM golang:1.26.4-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git ca-certificates tzdata

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest AS runner

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/server ./server

EXPOSE 8080

CMD ["./server"]
