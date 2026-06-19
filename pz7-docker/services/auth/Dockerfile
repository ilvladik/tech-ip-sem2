FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/auth

FROM alpine:3.19

RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /bin/app ./app

ENV AUTH_PORT=8087
ENV LISTEN_ADDRESS=0.0.0.0
EXPOSE 8087

CMD ["./app"]
