FROM golang:1.9
RUN mkdir -p /go/src/app
WORKDIR /go/src/
COPY . .
RUN go get github.com/streadway/amqp
RUN go build main.go 
RUN go build receive.go Ledger.go Payment.go
