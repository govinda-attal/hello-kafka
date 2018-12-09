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

	// For given request response scenarios, default broker to connect for reply
	pCfg := &kafka.ConfigMap{"bootstrap.servers": "broker"}

	// To consume/listen messages - consumer connects to given broker
	cCfg := &kafka.ConfigMap{
		"bootstrap.servers": "broker",
		"group.id":          "myGroup",
		"auto.offset.reset": "earliest",
	}

	// This is actual service devoid of any handlers (no marshal/unmarshal)
	greeter := &Greeter{}
	// This is handler/mux very similar to http mux
	h := NewGreeterHandler(greeter)

	// Serivce will listen to given topic and any un-handled runtime errors will be written to 'errTopic'
	r := krouter.New(cCfg, pCfg, "greetTopic", "errTopic")

	// This is like service name
	g := r.NewRouteGrp("Greetings")

	// This is like an operation or on a given message for example 'Hello' message on service 'Greetings'
	// If service gets message with ReplyTo topic (within message headers)
	// then service will reply the message on that topic
	g.Invoke("Hello", h.Hello)

	stop := make(chan interface{})

	// Listen to messages
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
