FROM golang:1.22-alpine as builder

WORKDIR /app

COPY . .

RUN go build -o main cmd/main.go

FROM alpine:3.16

WORKDIR /app

COPY --from=builder /app/main .
COPY templates ./templates

# Run the application
CMD ["./main"]