package worker

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	util "github.com/Mekawy5/binance/util"
	"github.com/gorilla/websocket"
)

const (
	apiBase    = "https://api.binance.com/api/v3"
	socketBase = "wss://stream.binance.com:9443"
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

// Trade info
type Trade struct {
	Symbol string `json:"s"`
	Time   int    `json:"T"`
	ID     int    `json:"t"`
	Price  string `json:"p"`
	Amount string `json:"q"`
}

// GetSymbols from the exchange server
func GetSymbols() Symbols {
	// use http.NewRequest to send headers(auth/encoded...etc)
	res, err := http.Get(apiBase + "/exchangeInfo")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	var syms Symbols

	json.NewDecoder(res.Body).Decode(&syms)

	return syms
}

// SubscriptionString Build subscription string for btc trades.
func SubscriptionString(s Symbols) string {
	strs := "/stream?streams="
	for _, sym := range s.Symbols {
		if sym.Base == "BTC" {
			str := strings.ToLower(sym.Symbol) + "@trade/"
			strs += str
		}
	}
	return util.TrimLastChar(strs)
}

func connect() *websocket.Conn {
	subStr := SubscriptionString(GetSymbols())
	c, _, err := websocket.DefaultDialer.Dial(socketBase+subStr, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println("socket connection without errors.")
	return c
	// remember to defer close this connection.
}

// NewProcessor func creates new processor var
func NewProcessor() *Processor {
	return &Processor{
		conn: connect(),
	}
}

func (p *Processor) handleMessages() {
	var t Trade
	for {
		_, msg, err := p.conn.ReadMessage()
		if err != nil {
			log.Fatal("Error reading from connection", err)
			return
		}

		err = json.Unmarshal([]byte(msg), &t)
		if err != nil {
			panic(err)
		}

		fmt.Println(string([]byte(msg)))
	}
}

// Process read authentication result
func (p *Processor) Process() {
	defer p.conn.Close()

	p.handleMessages()
}
