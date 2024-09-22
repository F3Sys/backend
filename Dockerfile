ARG GO_VERSION=1
FROM golang:${GO_VERSION}-alpine as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /app cmd/api/main.go

FROM alpine:3

COPY --from=builder /app /usr/local/bin/

CMD [ "app" ]