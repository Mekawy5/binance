package main

import "github.com/Mekawy5/binance/pkg/worker"

func main() {
	p := worker.NewProcessor()
	p.Process()

	// here create channel named trades & send it to process so all trades exists in this channel
	// create kafka producer
	// start handling producer response
	// start recieve messages from trades channel and write it to kafka using the producer's Produce()
}
