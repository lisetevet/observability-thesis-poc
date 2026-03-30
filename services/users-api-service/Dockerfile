FROM golang:1.25 as builder

WORKDIR /app/services/users-api-service

# Module files
COPY services/users-api-service/go.mod ./go.mod
COPY services/users-api-service/go.sum ./go.sum

# Shared observability module
COPY pkg /app/pkg

RUN go mod download

# Service source
COPY services/users-api-service ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/users-api ./cmd

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/users-api /app/users-api
COPY services/users-api-service/config/config.docker.yaml /app/config/config.docker.yaml

EXPOSE 8081
ENTRYPOINT ["/app/users-api"]