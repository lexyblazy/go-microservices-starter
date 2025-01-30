package vision

import (
	"fmt"
	"os"
	// "runtime"
	broker "lexyblazy.github.com/microservices-starter/pkg/broker"
	"time"
)

const SUB_TOPIC = "thors.foo"
const PUB_TOPIC = "vision.foo"

func Start() {
	fmt.Println("This is the vision microservice")

	natsUrl := os.Getenv("NATS_URL")
	b := broker.New(natsUrl)

	b.Publish(PUB_TOPIC, []byte("hello there from vision service"))

	b.Subscribe(SUB_TOPIC, func(data []byte) {
		fmt.Printf("Received message on %s topic, message = %s \n", SUB_TOPIC, string(data))

		time.Sleep(5 * time.Second)

		b.Publish(PUB_TOPIC, []byte("hello there from vision service"))

	})

}
