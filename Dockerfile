# ── Stage 1: Build ──────────────────────────────────────
FROM golang:1.22-alpine AS builder

WORKDIR /src
COPY go.mod main.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /ip-app-golang main.go

# ── Stage 2: Run ────────────────────────────────────────
FROM alpine:3.19

RUN adduser -D appuser
USER appuser

COPY --from=builder /ip-app-golang /ip-app-golang

EXPOSE 8080

ENTRYPOINT ["/ip-app-golang"]
