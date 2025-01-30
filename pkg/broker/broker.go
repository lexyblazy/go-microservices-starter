package broker

import (
	"log"
	"runtime"

	nats "github.com/nats-io/nats.go"
)

type Broker struct {
	nc *nats.Conn
}

type MessageHandler func(data []byte)

func (b *Broker) Publish(topic string, message []byte) {
	b.nc.Publish(topic, message)

	b.nc.Flush()

}

func (b *Broker) Subscribe(topic string, handler MessageHandler) {
	b.nc.Subscribe(topic, func(msg *nats.Msg) {
		handler(msg.Data)
	})

	b.nc.Flush()

	runtime.Goexit()

}

func (b *Broker) Close() {
	b.nc.Close()
}

func New(natsUrl string) *Broker {

	if len(natsUrl) < 1 {
		log.Fatal("NATS_URL cannot be empty string")
	}

	nc, err := nats.Connect(natsUrl)

	if err != nil {
		log.Fatal("Failed to connect nats")
	}

	return &Broker{nc}
}
