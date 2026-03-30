FROM golang:1.25.4-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN go build -o server ./cmd/server

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
ENTRYPOINT ["./server"]