FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/hezzltest ./cmd/

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/hezzltest .
COPY --from=builder /app/config.yml .

EXPOSE 3000
CMD ["./hezzltest"]