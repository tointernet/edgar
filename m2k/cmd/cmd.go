package cmd

import (
	"github.com/ralvescosta/gokit/mqtt"
	"github.com/tointernet/edgar/pkgs"
)

func NewCmd(container *pkgs.Container) error {
	container.MQTTDispatcher = mqtt.NewMQTTDispatcher(container.Logger, container.MQTTClient.Client())

	if err := NewMQTTConsumer(container); err != nil {
		return err
	}

	return nil
}
