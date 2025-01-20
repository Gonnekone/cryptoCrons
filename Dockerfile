FROM golang:alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o mainApp ./cmd/cryptoCrons/main.go
RUN go build -o mainMigrator ./cmd/migrator/migrator.go

FROM alpine:latest AS final
WORKDIR /app
COPY --from=builder /app/mainApp .
COPY --from=builder /app/mainMigrator .
COPY --from=builder /app/migrations ./migrations
EXPOSE 8082
RUN /bin/sh
CMD ["/bin/sh", "-c", "./mainMigrator --migrations-path=./migrations --direction=up && ./mainApp"]