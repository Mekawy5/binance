package main

import "github.com/Mekawy5/binance/pkg/worker"

func main() {
	p := worker.NewProcessor()
	p.Process()
}
