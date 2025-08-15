# Stage 1: Build the application
FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -mod=mod -buildvcs=false -ldflags="-w -s" -o ./bin/server ./cmd/app

# Stage 2: Create the final, minimal image
FROM alpine:3.20.1

WORKDIR /app

COPY --from=builder /app/bin/server .

COPY config ./config

EXPOSE 8080

CMD ["./server"] 