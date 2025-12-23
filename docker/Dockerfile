FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git ca-certificates build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main .

RUN mkdir -p /app/files

EXPOSE 8080

CMD ["./main"]
