# Preferred version : 1.16 Alpine - Low in size, but not compatible with consul 
FROM golang:latest


WORKDIR /app

COPY go.mod .
COPY go.sum .

COPY . .

EXPOSE 4200

RUN go build

CMD ["./hertz_server"]