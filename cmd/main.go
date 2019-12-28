package main

import (
	"github.com/Mekawy5/binance/pkg/kafka"
	"github.com/Mekawy5/binance/pkg/worker"
)

func main() {
	tc := make(chan string)

	p := worker.NewProcessor(tc)
	go p.Process()

	kp := kafka.NewProducer()
	go kp.ProcessResponse()

	for msg := range tc {
		kp.Produce(msg)
	}
}
