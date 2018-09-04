package main

import (
	"bytes"
	"fmt"
	"github.com/streadway/amqp"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func failOnError(err error, msg string) {
	if err != nil {
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}
func Exists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
func check(e error) {
	if e != nil {
		panic(e)
	}
}
func readLastLedger() string {

	if Exists("/data/lastLedger.txt") {
		dat, err := ioutil.ReadFile("/data/lastLedger.txt")
		if err != nil {
			// handle error
		}
		return strings.TrimSpace(string(dat))

	}
	return "19681277"
}
func writeLastLedgerNumber(s string) {
	//        fmt.Println("write", s)
	d1 := []byte(s)
	err := ioutil.WriteFile("/data/lastLedger.txt", d1, 0644)
	check(err)
}

func main() {

	//Establish RabbitMQ connection
	conn, err := amqp.Dial(os.Getenv("AMQPConn"))
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()
	fmt.Println("Opened amqp connection")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()
	fmt.Println("opened channel")

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
	fmt.Println("Exchanged declared\nname:stellar\ntype:direct\ndurable:true\nauto-deleted:false\ninternal:false\nno-wait:false\nargs:nil")

	for {
		lastLedgerNumber := readLastLedger()
		nextLedgerNumber, err := strconv.ParseInt(lastLedgerNumber, 10, 32)
		nextLedgerNumber = nextLedgerNumber + 1
		if err != nil {

		}
		url := fmt.Sprintf("https://horizon.stellar.org/ledgers/%v/payments?limit=200", nextLedgerNumber)

		body := bytes.NewReader([]byte(""))
		req, err := http.NewRequest("GET", url, body)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			// handle err
		}
		defer resp.Body.Close()
		content, _ := ioutil.ReadAll(resp.Body)
		if strings.Contains(string(content), "Resource Missing") {
			fmt.Println(fmt.Sprintf("Ledger not found number: %v", strconv.FormatInt(nextLedgerNumber, 10)))
			time.Sleep(1000 * time.Millisecond)
			continue
		}
		fmt.Println((strconv.FormatInt(nextLedgerNumber, 10)))
		writeLastLedgerNumber(strconv.FormatInt(nextLedgerNumber, 10))

		err = ch.Publish(
			"stellar",            // exchange
			"BlockchainListener", // routing key
			false,                // mandatory
			false,                // immediate
			amqp.Publishing{
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body:         []byte(content),
			})
		fmt.Println(" [x] Sent %s", string(content))
		failOnError(err, "Failed to publish a message")

	}
}
