FROM golang:latest
WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY . .

RUN go build

CMD ["./hertz_server"]