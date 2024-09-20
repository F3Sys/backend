ARG GO_VERSION=1
FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /usr/src/backend
COPY /backend/go.mod /backend/go.sum ./
RUN go mod download && go mod verify
COPY /backend/ .
RUN go build -v -o /backend cmd/api/main.go

RUN backend

FROM caddy:2-alpine

COPY Caddyfile /etc/caddy/Caddyfile

RUN caddy run
# COPY --from=builder /backend /usr/share/caddy/