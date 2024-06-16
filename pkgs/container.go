package pkgs

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/ralvescosta/gokit/configs"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/mqtt"
	tinyHTTP "github.com/ralvescosta/gokit/tiny_http"
)

type (
	Container struct {
		Cfgs           *configs.Configs
		Logger         logging.Logger
		Sig            chan os.Signal
		MQTTClient     mqtt.MQTTClient
		MQTTPublisher  mqtt.MQTTPublisher
		MQTTDispatcher mqtt.MQTTDispatcher
		TinyHTTPServer tinyHTTP.TinyServer
	}
)

func NewContainer(cfgs *configs.Configs) (*Container, error) {
	logger, err := logging.NewDefaultLogger(cfgs.AppConfigs)
	if err != nil {
		return nil, err
	}

	sig, tinyServer, err := provideTinyServer(cfgs, logger)
	if err != nil {
		return nil, err
	}

	mqttClient, mqttDispatcher, mqttPublisher, err := provideMQTTClient(cfgs, logger)
	if err != nil {
		return nil, err
	}

	return &Container{
		Cfgs:           cfgs,
		Logger:         logger,
		Sig:            sig,
		MQTTClient:     mqttClient,
		MQTTDispatcher: mqttDispatcher,
		MQTTPublisher:  mqttPublisher,
		TinyHTTPServer: tinyServer,
	}, nil
}

func provideTinyServer(cfgs *configs.Configs, logger logging.Logger) (chan os.Signal, tinyHTTP.TinyServer, error) {
	logger.Debug("creating signal channel...")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	logger.Debug("signal channel created!")

	logger.Debug("creating tiny http server...")

	tinyServer := tinyHTTP.NewTinyServer(cfgs.HTTPConfigs, logger).
		Sig(sig).
		Prometheus()

	logger.Debug("tiny http server created!")

	return sig, tinyServer, nil
}

func provideMQTTClient(cfgs *configs.Configs, logger logging.Logger) (mqtt.MQTTClient, mqtt.MQTTDispatcher, mqtt.MQTTPublisher, error) {
	logger.Debug("connecting to mqtt...")

	mqttClient := mqtt.NewMQTTClient(cfgs, logger)
	err := mqttClient.Connect()
	if err != nil {
		return nil, nil, nil, err
	}

	logger.Debug("mqtt connected!")

	dispatcher := mqtt.NewMQTTDispatcher(logger, mqttClient.Client())
	publisher := mqtt.NewMQTTPublisher(logger, mqttClient.Client())

	return mqttClient, dispatcher, publisher, err
}
