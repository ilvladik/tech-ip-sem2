FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/app ./cmd/tasks

FROM alpine:3.19

RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /bin/app ./app

ENV TASKS_PORT=8088
ENV AUTH_BASE_URL=http://auth:8087
ENV LISTEN_ADDRESS=0.0.0.0
EXPOSE 8088

CMD ["./app"]
