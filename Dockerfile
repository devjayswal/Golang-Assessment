# Multi-stage build for smaller image
FROM golang:1.22-alpine AS build
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/bin/server ./cmd/server

FROM alpine:3.19
WORKDIR /app
COPY --from=build /app/bin/server /app/server
COPY --from=build /app/.env /app/.env
EXPOSE 8080
CMD ["/app/server"]
