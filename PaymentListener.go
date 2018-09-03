package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

var (
	outfile, _ = os.Create("/data/stellarPayments.log") // update path for your needs
	l          = log.New(outfile, "", 0)
)

func main() {

	fmt.Println("Entering PaymentListner")
	conn, err := amqp.Dial(os.Getenv("AMQPConn"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"stellar", // name
		"direct",  // type
		true,      // durable
		false,     // auto-deleted
		false,     // internal
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare an exchange")

	q, err := ch.QueueDeclare(
		"client1", // name
		true,      // durable
		false,     // delete when usused
		true,      // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(
		q.Name,               // queue name
		"BlockchainListener", // routing key
		"stellar",            // exchange
		false,
		nil)
	failOnError(err, "Failed to bind a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var ledger Ledger
			json.Unmarshal(d.Body, &ledger)
			for i := range ledger.Embedded.Record {
				if ledger.Embedded.Record[i].AssetType == "native" {
					ledger.Embedded.Record[i].AssetCode = "XLM"
				} else {
					ledger.Embedded.Record[i].AssetCode = fmt.Sprintf("%v.%v", ledger.Embedded.Record[i].AssetCode, ledger.Embedded.Record[i].AssetIssuer)
				}
				if ledger.Embedded.Record[i].Type == "payment" {

					payment := Payment{
						Currency: ledger.Embedded.Record[i].AssetCode,
						Address:  ledger.Embedded.Record[i].To,
						Amount:   ledger.Embedded.Record[i].Amount,
						Hash:     ledger.Embedded.Record[i].Hash,
					}

					paymentJSON, err := json.Marshal(payment)
					if err != nil {
						fmt.Println(err)
						return
					}
					err = ch.Publish(
						"stellar",  // exchange
						"payments", // routing key
						false,      // mandatory
						false,      // immediate
						amqp.Publishing{
							DeliveryMode: amqp.Persistent,
							ContentType:  "text/plain",
							Body:         []byte(paymentJSON),
						})
					l.Printf(" [x] Sent %s", payment.String())
					failOnError(err, "Failed to publish a message")
				}

			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever

}
