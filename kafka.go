package main

import (
	_ "log"

	"github.com/Shopify/sarama"
)

type Message struct {
	Message  string `json:"message"`
	Date     string `json:"date"`
	Facility string `json:"facility"`
	Severity string `json:"severity"`
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
