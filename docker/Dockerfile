FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Убедитесь, что копируются все файлы, включая HTML шаблоны
# COPY . . уже копирует всё, но для ясности:
COPY internal/module/auth/pkg/template/ ./internal/module/auth/pkg/template/

RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/internal/module/auth/pkg/template/ ./internal/module/auth/pkg/template/

RUN mkdir -p /app/files

EXPOSE 8080

CMD ["./main", "-migrate=true"]
