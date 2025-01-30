package thor

import (
	"fmt"
	"os"
	"time"

	broker "lexyblazy.github.com/microservices-starter/pkg/broker"
)

const SUB_TOPIC = "vision.foo"
const PUB_TOPIC = "thors.foo"

func Start() {

	fmt.Println("this is the thor microservice")

	natsUrl := os.Getenv("NATS_URL")
	b := broker.New(natsUrl)

	b.Publish(PUB_TOPIC, []byte("hello there from thor service"))

	b.Subscribe(SUB_TOPIC, func(data []byte) {
		fmt.Printf("Received message on %s topic, message = %s \n", SUB_TOPIC, string(data))
		time.Sleep(2 * time.Second)
		b.Publish(PUB_TOPIC, []byte("hello there from thor service"))

	})

}
