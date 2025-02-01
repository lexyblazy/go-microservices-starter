package thor

import (
	"fmt"
	"time"

	broker "lexyblazy.github.com/microservices-starter/pkg/broker"
	"lexyblazy.github.com/microservices-starter/pkg/queue"
)

type Service struct {
	broker *broker.Broker
	queue  *queue.Queue
}

const SUB_TOPIC = "vision.foo"
const PUB_TOPIC = "thors.foo"
const STREAM_NAME = "thor_stream"

func New() *Service {
	broker := broker.New()
	queue := queue.New()

	return &Service{
		broker,
		queue,
	}
}

// test queue within the same service
func (s *Service) initQueue() {
	stream := s.queue.CreateStream(STREAM_NAME, []string{"thors_queue.*"})

	consumer := s.queue.CreateConsumer(stream)

	s.queue.Consume(consumer, func(data []byte) {
		fmt.Println("[Queue] => Received:", string(data))
	})

	counter := 1

	for {
		time.Sleep(2 * time.Second)
		s.queue.Publish("thors_queue.random", []byte(fmt.Sprintf("This is the %d message", counter)))
		counter += 1
	}
}

// test ping-pong between microservices
func (s *Service) initBroker() {
	s.broker.Publish(PUB_TOPIC, []byte("hello there from thor service"))

	s.broker.Subscribe(SUB_TOPIC, func(data []byte) {
		fmt.Printf("[PubSub] => Received message on %s topic, message = %s \n", SUB_TOPIC, string(data))
		time.Sleep(2 * time.Second)
		s.broker.Publish(PUB_TOPIC, []byte("hello there from thor service"))

	})

}

func Start() {

	fmt.Println("this is the thor microservice")
	done := make(chan bool)

	s := New()

	go s.initBroker()

	go s.initQueue()

	<-done
}
