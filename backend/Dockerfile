FROM golang:1.22-alpine AS build

WORKDIR /app
RUN apk add --no-cache ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /secscan ./main.go

FROM alpine:3.20

RUN apk add --no-cache ca-certificates wget && adduser -D -H appuser
WORKDIR /app

COPY --from=build /secscan /app/secscan

USER appuser
EXPOSE 8080
ENTRYPOINT ["/app/secscan"]
