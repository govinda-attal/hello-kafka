package main

import (
	"log"

	"github.com/govinda-attal/hello-kafka/pkg/example"
)

type Greeter struct{}

func (g Greeter) Hello(rq example.HelloRq) (example.HelloRs, error) {
	log.Println("message received from ", rq.Name)
	rs := example.HelloRs{Greetings: "Hello" + rq.Name}
	return rs, nil
}
