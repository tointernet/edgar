package pkg

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
	tinyHttp "github.com/ralvescosta/gokit/tiny_http"
)

type Container struct {
	cfgs               *configs.Configs
	logger             logging.Logger
	sig                chan os.Signal
	mqttClient         mqtt.MQTTClient
	prometheus         metrics.PrometheusMetrics
	prometheusShotdown func(context.Context) error
	tinyHTTPServer     tinyHttp.TinyServer
}

func NewContainer() (*Container, error) {
	cfgs, logger, err := provideConfigsAndLogger()
	if err != nil {
		return nil, err
	}

	prometheusProvider, shutdown, err := providePrometheus(cfgs, logger)
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
		cfgs:               cfgs,
		logger:             logger,
		sig:                sig,
		mqttClient:         mqttClient,
		prometheus:         prometheusProvider,
		prometheusShotdown: shutdown,
		tinyHTTPServer:     tinyServer,
	}, nil
}

func provideConfigsAndLogger() (*configs.Configs, logging.Logger, error) {
	cfgs, err := configsbuilder.NewConfigsBuilder().
		MQTT().
		Metrics().
		Build()

	if err != nil {
		return nil, nil, err
	}

	logger, err := logging.NewDefaultLogger(cfgs.AppConfigs)
	if err != nil {
		return nil, nil, err
	}

	return cfgs, logger, nil
}

func providePrometheus(cfgs *configs.Configs, logger logging.Logger) (metrics.PrometheusMetrics, func(context.Context) error, error) {
	logger.Debug("creating prometheus provider...")

	prometheusProvider := metrics.NewPrometheus(cfgs, logger)
	shutdown, err := prometheusProvider.Provider()
	if err != nil {
		return nil, nil, err
	}

	logger.Debug("prometheus provider created!")

	return prometheusProvider, shutdown, nil
}

func provideTinyServer(cfgs *configs.Configs, logger logging.Logger) (chan os.Signal, tinyHttp.TinyServer, error) {
	logger.Debug("creating signal channel...")

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	logger.Debug("signal channel created!")

	logger.Debug("creating tiny http server...")

	tinyServer := tinyHttp.NewTinyServer(cfgs.HTTPConfigs, logger).
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
