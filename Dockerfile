FROM golang:1.23-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod init rancher-ai-llm-mock || true
RUN go mod tidy
RUN go build -o llm-mock ./cmd/main.go

FROM alpine:3.19

WORKDIR /app
COPY --from=builder /app/llm-mock .
EXPOSE 8083
ENTRYPOINT ["./llm-mock"]
