package broker

import (
	"log"
	"runtime"

	nats "github.com/nats-io/nats.go"
	"lexyblazy.github.com/microservices-starter/pkg/common"
)

type Broker struct {
	nc *nats.Conn
}

func (b *Broker) Publish(topic string, message []byte) {
	b.nc.Publish(topic, message)

	b.nc.Flush()

}

func (b *Broker) Subscribe(topic string, handler common.MessageHandler) {
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

	common.LogFatalOnErr(err)

	b := &Broker{nc}

	go common.Teardown(b)

	return b
}
