FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod ./
RUN go mod download 
RUN go mod tidy 
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o collector ./cmd/api/main.go

FROM alpine:latest
COPY --from=builder /app/collector .
EXPOSE 8080
CMD ["./collector"]