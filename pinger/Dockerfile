FROM golang:1.23

WORKDIR /usr/src/pinger

COPY go.mod ./
RUN go mod download && go mod verify

COPY . .
RUN go build -v -o /usr/local/bin/pinger ./...

CMD ["pinger"]


