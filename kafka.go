package main

import (
	//"encoding/json"
	//	"fmt"
	"github.com/Shopify/sarama"
	_ "log"
	//"os"
	//"os/signal"
	//"regexp"
)

type Message struct {
	Message string `json:"message"`
}

type KafkaConsumer struct {
	broker   string
	topic    string
	consumer sarama.Consumer
}

//const (
//	prismPattern string = `(?P<accountId>[0-9]+)/(?P<appId>[0-9]+)`
//)

func NewKafkaConsumer(topic string, brokers []string) (KafkaConsumer, error) {
	k := KafkaConsumer{}
	//offset := sarama.OffsetOldest
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := sarama.NewConsumer(brokers, config)
	if err != nil {
		panic(err)
	}
	k.topic = topic
	k.consumer = consumer
	if err != nil {
		// return nil, err
		panic(err)
	}
	return k, nil
}

func (k *KafkaConsumer) close() error {
	if err := k.consumer.Close(); err != nil {
		panic(err)
	}
	return nil
}

/*
func (k *KafkaConsumer) read() error {
	consumer, err := k.consumer.ConsumePartition(k.topic, 0, sarama.OffsetOldest)
	if err != nil {
		panic(err)
	}

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	// Count message processed
	msgCount := 0
	// Get signnal for finish
	doneCh := make(chan struct{})
	r := regexp.MustCompile(prismPattern)
	go func() {
		for {
			select {
			case err := <-consumer.Errors():
				fmt.Println(err)
			case msg := <-consumer.Messages():
				msgCount++
				message := Message{}
				err := json.Unmarshal([]byte(msg.Value), &message)
				if err != nil {
					// TODO: handle unmarshal err here
					fmt.Println("unmarshal failed")
				}
				a := r.FindStringSubmatch(string(message.Message))
				fmt.Println(a)
				if len(a) != 0 {
					SentToFile([]byte(message.Message), string(a[1]))
				}
			case <-signals:
				fmt.Println("interrupt is detected")
				doneCh <- struct{}{}
			}
		}
	}()
	<-doneCh
	fmt.Println("Processed", msgCount, "messages")
	return nil
}
*/
