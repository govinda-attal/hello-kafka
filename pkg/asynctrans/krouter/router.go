package krouter

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"github.com/confluentinc/confluent-kafka-go/kafka"

	at "github.com/govinda-attal/hello-kafka/pkg/asynctrans"
	"github.com/govinda-attal/hello-kafka/pkg/core/status"
)

type router struct {
	consumerCfg, producerCfg *kafka.ConfigMap
	routeGrps                map[string]*RouteGroup
	rqTopic, errTopic        string
	NotFoundHandler          at.MsgHandler
}

func New(consumerCfg, producerCfg *kafka.ConfigMap, rqTopic, errTopic string) *router {
	return &router{
		consumerCfg: consumerCfg,
		producerCfg: producerCfg,
		rqTopic:     rqTopic,
		errTopic:    errTopic,
		routeGrps:   make(map[string]*RouteGroup),
	}
}

func (ks *router) NewRouteGrp(grpName string) *RouteGroup {
	rg := &RouteGroup{group: grpName, handlers: make(map[string]at.MsgHandler)}
	ks.routeGrps[grpName] = rg
	return rg
}

func (ks *router) Listen(stop chan interface{}) error {

	log.Println("Consumer is now in preparation!")

	c, err := kafka.NewConsumer(ks.consumerCfg)

	if err != nil {
		panic(err)
	}

	c.SubscribeTopics([]string{ks.rqTopic, "^aRegex.*[Tt]opic"}, nil)

	// OPTION - 1 // WORKS
	// for {
	// 	msg, err := c.ReadMessage(-1)
	// 	if err == nil {
	// 		ks.callHandler(msg)
	// 	} else {
	// 		fmt.Printf("Consumer error: %v (%v)\n", err, msg)
	// 	}
	// }

	// OPTION - 2 // NOT WORKS
	// for {
	// 	select {

	// 	case <-stop:
	// 		break
	// 	case ev := <-c.Events():
	// 		switch e := ev.(type) {
	// 		case *kafka.Message:
	// 			log.Println("Message received: ", e.Key, "\n", e.Value)
	// 			ks.callHandler(e)
	// 		default:
	// 			log.Println("Unknown event ", ev)
	// 		}
	// 	}
	// }

	// OPTION 3 - WORKS

	run := true

	for run == true {
		select {
		case <-stop:
			fmt.Println("Caught signal - terminating!")
			run = false
		default:
			ev := c.Poll(100)
			if ev == nil {
				continue
			}

			switch e := ev.(type) {
			case *kafka.Message:
				ks.callHandler(e)
			case kafka.Error:
				fmt.Println("Error: ", e)
				run = false
			default:
				fmt.Println("Ignored: ", e)
			}
		}
	}

	fmt.Println("Closing consumer!")
	c.Close()

	return nil
}

func (ks *router) callHandler(msg *kafka.Message) error {

	msgKey := msg.Key

	msgType := headerByKey(msg.Headers, at.MsgHdrMsgType)
	msgName := headerByKey(msg.Headers, at.MsgHdrMsgName)
	grpName := headerByKey(msg.Headers, at.MsgHdrGrpName)
	replyTo := headerByKey(msg.Headers, at.MsgHdrReplyTo)

	ctx := context.Background()
	ctx = context.WithValue(ctx, at.CtxKeyMsgID, string(msg.Key))

	h, err := ks.GetMsgHandler(grpName, msgName)

	if err != nil {
		return err
	}

	rs, err := h(ctx, msg.Value)

	if msgType == at.MsgTypeRq {
		err = ks.reply(msgName, msgKey, replyTo, rs, err)
	}

	if err != nil {
		ks.write(at.MsgTypeErrEvent, msgName, msgKey, ks.errTopic, nil, err)
	}

	return err
}

func (ks *router) GetMsgHandler(grpName, msgName string) (at.MsgHandler, error) {
	rg, ok := ks.routeGrps[grpName]
	if !ok {
		return nil, status.ErrBadRequest.WithMessage(fmt.Sprintf("handler for given group name '%s' not found", grpName))
	}
	return rg.GetMsgHandler(msgName)
}

func (ks *router) reply(msgName string, msgKey []byte, rsTopic string, rs []byte, err error) error {
	return ks.write(at.MsgTypeRs, msgName, msgKey, rsTopic, rs, err)
}

func (ks *router) write(msgType at.MsgType, msgName string, msgKey []byte, topic string, data []byte, err error) error {
	if err != nil {
		errSvc, ok := err.(status.ErrServiceStatus)
		if !ok {
			errSvc = status.ErrInternal.WithError(err)
		}

		data, _ = json.Marshal(errSvc)

		msgType = at.MsgTypeErrEvent

		if msgType == at.MsgTypeRs {
			msgType = at.MsgTypeErrRs
		}
	}

	p, err := kafka.NewProducer(ks.producerCfg)
	if err != nil {
		panic(err)
	}

	defer p.Close()

	return p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Key:   msgKey,
		Value: data,
		Headers: []kafka.Header{
			kafka.Header{Key: at.MsgHdrMsgName, Value: []byte(msgName)},
			kafka.Header{Key: at.MsgHdrMsgType, Value: []byte(msgType)},
		},
	}, nil)
}

func headerByKey(hdrs []kafka.Header, key string) at.MsgType {
	for _, h := range hdrs {
		if h.Key == key {
			return string(h.Value)
		}
	}
	return at.MsgHdrValUnk
}
