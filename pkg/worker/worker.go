package worker

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/websocket"
)

// Processor struct contains websocket connection
type Processor struct {
	conn *websocket.Conn
}

// Symbol definition
type Symbol struct {
	Symbol string `json:"symbol"`
	Base   string `json:"baseAsset"`
	Quote  string `json:"quoteAsset"`
}

// Symbols response
type Symbols struct {
	Symbols []Symbol `json:"symbols"`
}

// GetSymbols from the exchange server
func GetSymbols() Symbols {
	// use http.NewRequest to send headers(auth/encoded...etc)
	res, err := http.Get("https://api.binance.com/api/v3/exchangeInfo")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var syms Symbols

	json.NewDecoder(res.Body).Decode(&syms)

	return syms
}
