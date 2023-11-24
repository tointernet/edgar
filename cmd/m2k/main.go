package main

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

func main() {
	conf := kafka.ConfigMap{}
	conf.SetKey("bootstrap.servers", "localhost:9092")

	topic := "test"
	p, err := kafka.NewProducer(&conf)
	if err != nil {
		panic(err)
	}

	p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte{},
		Headers: []kafka.Header{
			{Key: "myCustomHeader", Value: []byte("header value")},
		},
	}, nil)

}
