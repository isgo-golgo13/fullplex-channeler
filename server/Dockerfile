# Stage 1: Build
FROM golang:1.20-alpine AS builder
WORKDIR /app
COPY .. .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server ./server

# Stage 2: Run
FROM alpine:3.17
RUN addgroup -S appgroup && adduser -S appuser -G appgroup
WORKDIR /app
COPY --from=builder /app/server .
USER appuser
CMD ["./server"]
