FROM golang:1.19

WORKDIR /app
COPY . .

RUN go build -o proxy .

CMD ["./proxy"]
