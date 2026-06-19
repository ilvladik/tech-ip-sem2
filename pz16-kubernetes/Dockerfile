FROM golang:1.22-alpine AS builder
WORKDIR /src
COPY go.mod ./
COPY cmd/ ./cmd/
RUN CGO_ENABLED=0 GOOS=linux go build -o /tasks ./cmd/tasks

FROM alpine:3.20
WORKDIR /app
COPY --from=builder /tasks /app/tasks
EXPOSE 8082
CMD ["/app/tasks"]
