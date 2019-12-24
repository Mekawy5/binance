package worker

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	jsoniter "github.com/json-iterator/go"

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
	Symbol string
	Time   int
	ID     int
	Price  string
	Amount string
}

var json = jsoniter.ConfigCompatibleWithStandardLibrary

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
	for {
		_, msg, err := p.conn.ReadMessage()
		if err != nil {
			log.Fatal("Error reading from connection", err)
			return
		}

		data := json.Get(msg, "data").ToString()
		if data != "" {
			t := Trade{
				Symbol: json.Get(msg, "data", "s").ToString(),
				ID:     json.Get(msg, "data", "t").ToInt(),
				Time:   json.Get(msg, "data", "T").ToInt(),
				Price:  json.Get(msg, "data", "p").ToString(),
				Amount: json.Get(msg, "data", "q").ToString(),
			}

			msg, _ := json.MarshalToString(t)

			fmt.Println(msg)
			// Produce message to kafka.
		}
	}
}

// Process read authentication result
func (p *Processor) Process() {
	defer p.conn.Close()

	p.handleMessages()
}
