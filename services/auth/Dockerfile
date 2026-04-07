FROM golang:1.25.5-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /auth ./services/auth/cmd/auth

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /auth .

EXPOSE 8081 50051

ENV HTTP_PORT=8081
ENV GRPC_PORT=50051

CMD ["./auth"]
