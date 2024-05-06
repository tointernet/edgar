package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	topic    = "topic_0"
	server   = ""
	userName = ""
	password = ""
)

func main() {
	forever := make(chan os.Signal, 1)
	signal.Notify(forever, syscall.SIGINT, syscall.SIGTERM)

	consumer()
	// producer()

	<-forever
	log.Println("RECEIVED SIGTERM, FINISHING THE APPLICATION")
}

func consumer() {
	log.Println("[consumer] connection to kafka...")

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  server,
		"group.id":           "basic-consumer",
		"session.timeout.ms": 6000,
		"security.protocol":  "SASL_SSL",
		"sasl.mechanisms":    "PLAIN",
		"sasl.username":      userName,
		"sasl.password":      password,
		"enable.auto.commit": false,
	})

	if err != nil {
		log.Fatal("[consumer] connection failure", err)
	}

	log.Println("[consumer] connected...")

	log.Println("[consumer] subscribing to a topic...")
	consumer.SubscribeTopics([]string{topic}, nil)
	log.Println("[consumer] subscripted!")

	go func() {
		for {
			ev := consumer.Poll(100)
			if ev == nil {
				fmt.Print(".")
				continue
			}

			fmt.Print("\n")
			switch e := ev.(type) {
			case *kafka.Message:
				log.Println("[consumer] received message!")
				log.Printf("[consumer] topic: \"%v\", value: \"%v\"\n", *e.TopicPartition.Topic, string(e.Value))
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
			default:
				log.Printf("[consumer] Ignored event - %v\n", e)
			}
		}
	}()
}

func producer() {
	log.Println("[publisher] connection to kafka...")

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": server,
		"compression.codec": "gzip",
		"compression.type":  "gzip",
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     userName,
		"sasl.password":     password,
		"retries":           3,
		"retry.backoff.ms":  18000,
	})
	if err != nil {
		log.Fatal("[publisher] connection failure", err)
	}

	log.Println("[publisher] connected...")

	log.Println("[publisher] publishing data...")

	deliveryChannel := make(chan kafka.Event, 10000)
	err = p.Produce(&kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: kafka.PartitionAny,
		},
		Value: []byte("oi eu sou o goku"),
		Headers: []kafka.Header{
			{Key: "myCustomHeader", Value: []byte("header value")},
		},
		Key:           []byte("eventType"),
		Timestamp:     time.Now(),
		TimestampType: kafka.TimestampCreateTime,
	},
		deliveryChannel,
	)

	if err != nil {
		log.Fatal("[publisher] publish failure", err)
	}

	log.Println("[publisher] data was published")
}
