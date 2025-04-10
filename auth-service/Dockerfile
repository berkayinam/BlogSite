FROM golang:1.23.4 as builder

WORKDIR /app
COPY . .
RUN go build -o auth-service .

FROM alpine:latest
RUN apk add --no-cache libc6-compat
WORKDIR /
COPY --from=builder /app/auth-service ./
ENV JWT_SECRET=supersecretjwtkey
RUN chmod +x /auth-service
EXPOSE 8085
CMD ["./auth-service"]

