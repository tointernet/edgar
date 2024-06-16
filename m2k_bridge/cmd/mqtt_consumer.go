package cmd

import (
	"time"

	"github.com/ralvescosta/gokit/mqtt"
	"github.com/tointernet/edgar/pkgs"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"go.uber.org/zap"
)

func NewMQTTConsumer(container *pkgs.Container) error {
	container.MQTTDispatcher = mqtt.NewMQTTDispatcher(container.Logger, container.MQTTClient.Client())

	container.Logger.Debug("[publisher] connection to kafka...")

	server, username, password, topic := "", "", "", ""

	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": server,
		"compression.codec": "gzip",
		"compression.type":  "gzip",
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"sasl.username":     username,
		"sasl.password":     password,
		"retries":           3,
		"retry.backoff.ms":  18000,
	})
	if err != nil {
		container.Logger.Fatal("[publisher] connection failure", zap.Error(err))
	}

	container.Logger.Debug("[publisher] connected...")

	container.Logger.Debug("[publisher] publishing data...")

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
		container.Logger.Fatal("[publisher] publish failure", zap.Error(err))
	}

	container.Logger.Debug("[publisher] data was published")

	return nil
}
