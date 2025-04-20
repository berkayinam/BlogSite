FROM golang:1.21-alpine as builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download && go mod tidy

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o main ./cmd/main.go

FROM alpine:latest

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Expose the port
EXPOSE 8084

# Run the application
CMD ["./main"] 