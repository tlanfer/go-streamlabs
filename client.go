package streamlabs

import (
	"fmt"
	gosocketio "github.com/ambelovsky/gosf-socketio"
	"github.com/ambelovsky/gosf-socketio/transport"
	"log"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	Token string
}

type Donation struct {
	Amount   int
	Currency string
}

func (s *Client) Connect(events chan<- Donation) error {
	websocketTransport := transport.GetDefaultWebsocketTransport()
	websocketTransport.PingInterval = 5 * time.Second

	client, err := gosocketio.Dial(gosocketio.GetUrl("sockets.streamlabs.com", 443, true)+"&token="+s.Token, websocketTransport)
	if err != nil {
		return fmt.Errorf("failed to subscribe to create client: %w", err)
	}

	err = client.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		log.Printf("Streamlabs connected")
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to connects: %w", err)
	}

	err = client.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		log.Printf("Streamlabs disconnected")
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to disconnects: %w", err)
	}

	err = client.On("event", func(c *gosocketio.Channel, data Ev) {
		if data.Type == "donation" {
			amount := parseAmount(data.Message[0].Amount)
			currency := strings.ToLower(data.Message[0].Currency)

			events <- Donation{
				Amount:   amount,
				Currency: currency,
			}
		}
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to events: %w", err)
	}

	return nil
}

func parseAmount(input interface{}) int {
	asString := fmt.Sprint(input)
	num, err := strconv.ParseFloat(asString, 32)
	if err != nil {
		return -1
	}
	return int(num * 100)
}

type Ev struct {
	For     string `json:"for"`
	Type    string `json:"type"`
	Message []struct {
		Amount   interface{} `json:"amount"`
		Currency string      `json:"currency"`
	} `json:"message"`
}
