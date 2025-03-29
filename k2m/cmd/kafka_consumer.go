package cmd

import (
	"fmt"
	"os"

	"github.com/ralvescosta/gokit/mqtt"
	"github.com/tointernet/edgar/pkgs"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

func NewKafkaConsumer(container *pkgs.Container) error {
	container.MQTTPublisher = mqtt.NewMQTTPublisher(container.Logger, container.MQTTClient.Client())

	container.Logger.Debug("[consumer] connection to kafka...")

	server, username, password, topic := "", "", "", ""

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  server,
		"group.id":           "basic-consumer",
		"session.timeout.ms": 6000,
		"security.protocol":  "SASL_SSL",
		"sasl.mechanisms":    "PLAIN",
		"sasl.username":      username,
		"sasl.password":      password,
		"enable.auto.commit": false,
	})

	if err != nil {
		container.Logger.Fatal("[consumer] connection failure", zap.Error(err))
	}

	container.Logger.Debug("[consumer] connected...")

	container.Logger.Debug("[consumer] subscribing to a topic...")
	consumer.SubscribeTopics([]string{topic}, nil)
	container.Logger.Debug("[consumer] subscripted!")

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
				container.Logger.Debug("[consumer] received message!")
				container.Logger.Debug(fmt.Sprintf("[consumer] topic: \"%v\", value: \"%v\"\n", *e.TopicPartition.Topic, string(e.Value)))
			case kafka.Error:
				// Errors should generally be considered
				// informational, the client will try to
				// automatically recover.
				fmt.Fprintf(os.Stderr, "%% Error: %v: %v\n", e.Code(), e)
			default:
				container.Logger.Debug(fmt.Sprintf("[consumer] Ignored event - %v\n", e))
			}
		}
	}()
	return nil
}
