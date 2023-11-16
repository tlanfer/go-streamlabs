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

type Streamlabs interface {
	Connect(token string) error
	OnDonation(c chan<- Donation)
}

func New() Streamlabs {
	return &client{}

}

type client struct {
	chDonations chan<- Donation
	client      *gosocketio.Client
}

var (
	Log = log.Default()
)

func (s *client) Connect(token string) error {
	if s.client != nil {
		s.client.Close()
	}

	websocketTransport := transport.GetDefaultWebsocketTransport()
	websocketTransport.PingInterval = 5 * time.Second

	client, err := gosocketio.Dial(gosocketio.GetUrl("sockets.streamlabs.com", 443, true)+"&token="+token, websocketTransport)
	if err != nil {
		return fmt.Errorf("failed to subscribe to create client: %w", err)
	}

	s.client = client

	err = client.On(gosocketio.OnConnection, func(c *gosocketio.Channel) {
		Log.Printf("Streamlabs connected")
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to connects: %w", err)
	}

	err = client.On(gosocketio.OnDisconnection, func(c *gosocketio.Channel) {
		Log.Printf("Streamlabs disconnected")
	})
	if err != nil {
		return fmt.Errorf("failed to subscribe to disconnects: %w", err)
	}

	err = client.On("event", func(c *gosocketio.Channel, data ev) {
		switch data.Type {
		case "donation":
			amount := parseAmount(data.Message[0].Amount) * 100
			currency := strings.ToLower(data.Message[0].Currency)

			if s.chDonations != nil {
				s.chDonations <- Donation{
					Amount:   amount,
					Currency: currency,
				}
			}
		}
	})

	if err != nil {
		return fmt.Errorf("failed to subscribe to events: %w", err)
	}

	return nil
}

func (s *client) OnDonation(c chan<- Donation) {
	s.chDonations = c
}

func parseAmount(input interface{}) int {
	asString := fmt.Sprint(input)
	num, err := strconv.ParseFloat(asString, 32)
	if err != nil {
		return -1
	}
	return int(num)
}
