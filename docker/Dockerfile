FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download -x

COPY . .

WORKDIR /app/cmd/
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 go build -ldflags="-w -s" -o hezzltest .

FROM alpine:3.19

WORKDIR /app

COPY --from=builder /app/cmd/hezzltest .

EXPOSE 3000

CMD ["./hezzltest"]