FROM golang:latest

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -o /docker-f3s-backend

EXPOSE 8080

CMD ["/docker-f3s-backend"]
