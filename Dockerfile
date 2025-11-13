FROM golang:1.24.0-alpine AS builder

RUN apk add --no-cache git
RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o /usr/local/bin/avito-pr ./cmd/pr

FROM alpine:latest

ENV PORT=8080

COPY --from=builder /usr/local/bin/avito-pr /usr/local/bin/avito-pr
COPY --from=builder /go/bin/migrate /usr/local/bin/migrate

WORKDIR /app

EXPOSE ${PORT}

CMD ["avito-pr"]