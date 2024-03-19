# 1st stage
FROM golang:1.22-alpine3.19 AS builder

WORKDIR /compiler

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/main.go

# 2nd stage
FROM alpine:3 AS prod

WORKDIR /app

COPY --from=builder /compiler/server .

CMD ["./server"]