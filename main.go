package main

import (
	"bytes"
	"fmt"
	//	"github.com/streadway/amqp"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
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

func readLastBlock() string {

	if Exists("lastLedger.txt") {
		dat, err := ioutil.ReadFile("lastLedger.txt")
		if err != nil {
			// handle error
		}
		return strings.TrimSpace(string(dat))

	}
	return "19681277"
}

func main() {
	//Establish RabbitMQ connection
	/*	conn, err := amqp.Dial(os.Getenv("AMQPConn"))
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()
		log.Printf("Opened amqp connection")

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()
		log.Printf("opened channel")

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
		log.Printf("Exchanged declared\nname:stellar\ntype:direct\ndurable:true\nauto-deleted:false\ninternal:false\nno-wait:false\nargs:nil")

		//	for {
	*/
	lastBlockNumber := readLastBlock()
	nextBlockNumber, err := strconv.ParseInt(lastBlockNumber, 10, 32)
	if err != nil {

	}
	url := fmt.Sprintf("https://horizon.stellar.org/ledgers/%v/payments?limit=200", nextBlockNumber)

	body := bytes.NewReader([]byte(""))
	req, err := http.NewRequest("GET", url, body)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle err
	}
	defer resp.Body.Close()
	content, _ := ioutil.ReadAll(resp.Body)
	var l Ledger
	json.Unmarshal(content, &l)
	fmt.Println(l.Embedded.Record[1])
	//}

}
