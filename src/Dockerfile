FROM golang:1.26-alpine AS builder

ARG VERSION=0.0.0
ARG COMMIT_HASH=n/a
ARG BUILD_TIMESTAMP=n/a
ARG BUILD_TYPE=prod

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -ldflags="-X 'talents-api/app.Version=${VERSION}' -X 'talents-api/app.CommitHash=${COMMIT_HASH}' -X 'talents-api/app.BuildTimestamp=${BUILD_TIMESTAMP}' -X 'talents-api/app.Build=${BUILD_TYPE}'" -o main .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main .
COPY --from=builder /app/public ./public

EXPOSE 8080

CMD ["./main"]
