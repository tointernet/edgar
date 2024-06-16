package pkgs

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/ralvescosta/gokit/configs"
	configsbuilder "github.com/ralvescosta/gokit/configs_builder"
	"github.com/ralvescosta/gokit/logging"
	"github.com/ralvescosta/gokit/metrics"
	"github.com/ralvescosta/gokit/mqtt"
	tinyHTTP "github.com/ralvescosta/gokit/tiny_http"
)

type (
	Container struct {
		Cfgs               *configs.Configs
		Logger             logging.Logger
		Sig                chan os.Signal
		MQTTClient         mqtt.MQTTClient
		MQTTPublisher      mqtt.MQTTPublisher
		MQTTDispatcher     mqtt.MQTTDispatcher
		Prometheus         metrics.PrometheusMetrics
		PrometheusShotdown func(context.Context) error
		TinyHTTPServer     tinyHTTP.TinyServer
	}
)

func NewContainer() (*Container, error) {
	cfgs, err := configsbuilder.NewConfigsBuilder().
		MQTT().
		HTTP().
		Metrics().
		Build()

	if err != nil {
		return nil, err
	}

	logger, err := logging.NewDefaultLogger(cfgs.AppConfigs)
	if err != nil {
		return nil, err
	}

	sig, tinyServer, err := provideTinyServer(cfgs, logger)
	if err != nil {
		return nil, err
	}

	mqttClient, err := provideMQTTClient(cfgs, logger)
	if err != nil {
		return nil, err
	}

	return &Container{
		Cfgs:           cfgs,
		Logger:         logger,
		Sig:            sig,
		MQTTClient:     mqttClient,
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

func provideMQTTClient(cfgs *configs.Configs, logger logging.Logger) (mqtt.MQTTClient, error) {
	logger.Debug("connecting to mqtt...")

	mqttClient := mqtt.NewMQTTClient(cfgs, logger)
	err := mqttClient.Connect()
	if err != nil {
		return nil, err
	}

	logger.Debug("mqtt connected!")

	return mqttClient, err
}
