package main

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	at "github.com/govinda-attal/hello-kafka/pkg/asynctrans"
)

func main() {

	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "broker"})
	if err != nil {
		panic(err)
	}

	defer p.Close()

	// Delivery report handler for produced messages
	go func() {
		for e := range p.Events() {
			switch ev := e.(type) {
			case *kafka.Message:
				if ev.TopicPartition.Error != nil {
					fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
				} else {
					fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
				}
			}
		}
	}()

	// Produce messages to topic (asynchronously)
	topic := "myTopic"
	for _, _ = range []string{"", ""} {
		err := p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          []byte(`{"name":"govinda"}`),
			Headers:        []kafka.Header{kafka.Header{Key: at.MsgHdrGrpName, Value: []byte("Greetings")}, kafka.Header{Key: at.MsgHdrMsgType, Value: []byte(at.MsgTypeRq)}, kafka.Header{Key: at.MsgHdrMsgName, Value: []byte("Hello")}, kafka.Header{Key: at.MsgHdrReplyTo, Value: []byte("replyTopic")}},
		}, nil)
		if err != nil {
			fmt.Println(err)
		}
	}

	// Wait for message deliveries before shutting down
	p.Flush(15 * 1000)
}
