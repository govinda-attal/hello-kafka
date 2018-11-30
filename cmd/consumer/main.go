// package main

// import (
// 	"fmt"

// 	"github.com/confluentinc/confluent-kafka-go/kafka"
// )

// func main() {

// 	c, err := kafka.NewConsumer(&kafka.ConfigMap{
// 		"bootstrap.servers": "broker",
// 		"group.id":          "myGroup",
// 		"auto.offset.reset": "earliest",
// 	})

// 	if err != nil {
// 		panic(err)
// 	}

// 	c.SubscribeTopics([]string{"myTopic", "^aRegex.*[Tt]opic"}, nil)

// 	for {
// 		msg, err := c.ReadMessage(-1)
// 		if err == nil {
// 			fmt.Printf("Message on %s: %s\n", msg.TopicPartition, string(msg.Value))
// 		} else {
// 			// The client will automatically try to recover from all errors.
// 			fmt.Printf("Consumer error: %v (%v)\n", err, msg)
// 		}
// 	}

// 	c.Close()
// }

package main

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	"github.com/govinda-attal/hello-kafka/pkg/asynctrans/krouter"
)

func main() {

	pCfg := &kafka.ConfigMap{"bootstrap.servers": "broker"}
	cCfg := &kafka.ConfigMap{
		"bootstrap.servers":               "broker",
		"group.id":                        "myGroup",
		"session.timeout.ms":              6000,
		"go.events.channel.enable":        true,
		"go.application.rebalance.enable": true,
		"default.topic.config":            kafka.ConfigMap{"auto.offset.reset": "earliest"},
	}

	greeter := &Greeter{}
	h := NewGreeterHandler(greeter)

	r := krouter.New(cCfg, pCfg, "myTopic", "errTopic")

	g := r.NewRouteGrp("Greetings")
	g.Invoke("Hello", h.Hello)

	stop := make(chan interface{})

	go func() {
		if err := r.Listen(stop); err != nil {
			log.Fatalln(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	stop <- struct{}{}

	time.Sleep(5 * time.Second)
}
