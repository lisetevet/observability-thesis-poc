FROM golang:1.25 as builder

WORKDIR /app/services/mobile-api-service

COPY services/mobile-api-service/go.mod ./go.mod
COPY services/mobile-api-service/go.sum ./go.sum

RUN go mod download

COPY services/mobile-api-service ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/mobile-api ./cmd

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/mobile-api /app/mobile-api
COPY services/mobile-api-service/config/config.docker.yaml /app/config/config.docker.yaml

EXPOSE 8082
ENTRYPOINT ["/app/mobile-api"]