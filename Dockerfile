FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /sisyphus ./cmd/sisyphus

FROM alpine:3.22

WORKDIR /app

COPY --from=builder /sisyphus /app/sisyphus

ENTRYPOINT ["/app/sisyphus"]
