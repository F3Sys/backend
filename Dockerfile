ARG GO_VERSION=1
FROM golang:${GO_VERSION}-bookworm as builder

WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /app cmd/api/main.go


FROM debian:bookworm
ARG LITEFS_CONFIG=litefs.yml

ADD etc/litefs.yml /tmp/litefs.yml
RUN cp /tmp/$LITEFS_CONFIG /etc/litefs.yml

COPY --from=flyio/litefs:0.5 /usr/local/bin/litefs /usr/local/bin/litefs
RUN apt-get update -y && apt-get install -y ca-certificates fuse3 sqlite3

COPY --from=builder /app /usr/local/bin/

ENTRYPOINT litefs mount
