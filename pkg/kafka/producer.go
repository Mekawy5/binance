package kafka

import (
	"fmt"

	"github.com/Shopify/sarama"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	brokerList = kingpin.Flag("brokerList", "List of brokers to connect").Default("localhost:9093").Strings()
	topic      = kingpin.Flag("topic", "Topic name").Default("trades").String()
	maxRetry   = kingpin.Flag("maxRetry", "Retry limit").Default("5").Int()
)

//Producer object holds async producer and topic
type Producer struct {
	t string
	sarama.AsyncProducer
}

// NewProducer creates producer object
func NewProducer() *Producer {
	kingpin.Parse()
	conf := sarama.NewConfig()
	conf.Producer.RequiredAcks = sarama.WaitForAll
	conf.Producer.Retry.Max = *maxRetry
	conf.Producer.Return.Successes = true
	conf.Producer.Return.Errors = true

	prod, err := sarama.NewAsyncProducer(*brokerList, conf)
	if err != nil {
		panic(err)
	}

	return &Producer{
		*topic,
		prod,
	}
}

// Produce message to defined topic, push message to producer input channel asynchronously.
func (p *Producer) Produce(msg string) {
	m := sarama.ProducerMessage{Topic: p.t, Value: sarama.StringEncoder(msg)}
	p.Input() <- &m
}

// ProcessResponse listens to responses channels to handle faliyre & success channels messages
func (p *Producer) ProcessResponse() {
	for {
		select {
		case result := <-p.Successes():
			fmt.Printf("Wrote to Partition : %d - offset : %d\n", result.Partition, result.Offset)
		case err := <-p.Errors():
			fmt.Println("Failed to produce message", err)
		}
	}
}
