package queue

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"lexyblazy.github.com/microservices-starter/pkg/common"
)

type Queue struct {
	nc            *nats.Conn
	js            jetstream.JetStream
	ctx           context.Context
	ctxCancel     context.CancelFunc
	unsubConsumer func()
}

func (q *Queue) CreateStream(streamName string, subjects []string) jetstream.Stream {

	s, err := q.js.CreateStream(q.ctx, jetstream.StreamConfig{
		Name:     streamName,
		Subjects: subjects,
	})

	common.LogFatalOnErr(err)

	return s
}

func (q *Queue) CreateConsumer(s jetstream.Stream) jetstream.Consumer {
	cons, err := s.CreateOrUpdateConsumer(q.ctx, jetstream.ConsumerConfig{
		AckPolicy: jetstream.AckExplicitPolicy,
	})

	common.LogFatalOnErr(err)

	return cons
}

func (q *Queue) Consume(consumer jetstream.Consumer, handler common.MessageHandler) {
	cc, err := consumer.Consume(func(msg jetstream.Msg) {
		handler(msg.Data())
		msg.Ack()
	}, jetstream.ConsumeErrHandler(func(consumeCtx jetstream.ConsumeContext, err error) {
		log.Println(err)
	}))

	common.LogFatalOnErr(err)

	// attach unsubscribe and cancel handler to Queue struct
	q.unsubConsumer = cc.Stop

}

func (q *Queue) Publish(topic string, message []byte) {

	if q.nc.Status() != nats.CONNECTED {
		log.Println("NATS is disconnected")

		return
	}

	if _, err := q.js.Publish(q.ctx, topic, message); err != nil {
		log.Println("pub error: ", err)
	}
}

func (q *Queue) Close() {
	q.ctxCancel()

	if q.unsubConsumer != nil {
		q.unsubConsumer()
	}

	q.nc.Close()
}

func New() *Queue {

	natsUrl := os.Getenv("NATS_URL")

	if len(natsUrl) < 1 {
		log.Fatal("NATS_URL cannot be empty string")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Minute)

	nc, err := nats.Connect(natsUrl)
	common.LogFatalOnErr(err)

	js, err := jetstream.New(nc)

	common.LogFatalOnErr(err)

	q := &Queue{
		nc: nc, js: js, ctx: ctx, ctxCancel: cancel,
		unsubConsumer: nil,
	}

	go common.Teardown(q)

	return q

}
