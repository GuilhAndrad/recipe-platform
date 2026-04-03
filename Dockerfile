# ─── Stage 1: Build ───────────────────────────────────────────────────────────
FROM golang:1.25.0-alpine AS builder

# gcc é necessário para o driver postgres (CGO)
RUN apk add --no-cache gcc musl-dev

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .

# Compila apenas o binário da API
RUN CGO_ENABLED=1 GOOS=linux go build -o /app/api ./cmd/api

# ─── Stage 2: Runtime ─────────────────────────────────────────────────────────
FROM alpine:3.19

RUN apk add --no-cache ca-certificates

WORKDIR /app

# Copia apenas o binário compilado — imagem final bem menor
COPY --from=builder /app/api .

EXPOSE 8080

CMD ["./api"]