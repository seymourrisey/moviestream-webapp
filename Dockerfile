FROM golang:1.25.4-alpine AS builder
WORKDIR /app

COPY Server/appServer/go.mod Server/appServer/go.sum ./
RUN go mod download

COPY Server/appServer .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/app /app/app
EXPOSE 8080
CMD ["/app/app"]
