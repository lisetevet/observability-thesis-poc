FROM golang:1.25 as builder

WORKDIR /app/services/profile-service

COPY services/profile-service/go.mod ./go.mod
COPY services/profile-service/go.sum ./go.sum

COPY pkg /app/pkg

RUN go mod download

COPY services/profile-service ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app/profile-service ./cmd

FROM gcr.io/distroless/base-debian12

WORKDIR /app
COPY --from=builder /app/profile-service /app/profile-service
COPY services/profile-service/config/config.docker.yaml /app/config/config.docker.yaml

EXPOSE 8080
ENTRYPOINT ["/app/profile-service"]