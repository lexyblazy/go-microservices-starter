package vision

import (
	"fmt"
	"time"

	broker "lexyblazy.github.com/microservices-starter/pkg/broker"
	"lexyblazy.github.com/microservices-starter/pkg/queue"
)

const SUB_TOPIC = "thors.foo"
const PUB_TOPIC = "vision.foo"
const VISION_STREAM = "vision_stream"

type Service struct {
	broker *broker.Broker
	queue  *queue.Queue
}

func New() *Service {

	b := broker.New()
	q := queue.New()

	return &Service{
		broker: b,
		queue:  q,
	}
}

func (s *Service) initQueue() {
	stream := s.queue.CreateStream(VISION_STREAM, []string{"vision_queue.*"})

	c := s.queue.CreateConsumer(stream)

	s.queue.Consume(c, func(data []byte) {
		fmt.Println("[Queue] => Received:", string(data))
	})

	counter := 1

	for {
		time.Sleep(2 * time.Second)
		s.queue.Publish("vision_queue.random", []byte(fmt.Sprintf("This is the %d message", counter)))
		counter += 1
	}
}

func (s *Service) initBroker() {
	s.broker.Publish(PUB_TOPIC, []byte("hello there from vision service"))

	s.broker.Subscribe(SUB_TOPIC, func(data []byte) {
		fmt.Printf("[PubSub] => Received message on %s topic, message = %s \n", SUB_TOPIC, string(data))

		time.Sleep(5 * time.Second)

		s.broker.Publish(PUB_TOPIC, []byte("hello there from vision service"))

	})
}

func Start() {
	fmt.Println("This is the vision microservice")

	done := make(chan bool)

	s := New()

	go s.initBroker()

	go s.initQueue()

	<-done

}
