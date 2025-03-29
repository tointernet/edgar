package main

import (
	"github.com/tointernet/edgar/pkgs"
	"go.uber.org/zap"

	"github.com/tointernet/edgar/m2k/cmd"
)

func main() {
	container, err := pkgs.NewContainer()
	if err != nil {
		panic(err)
	}

	if err := cmd.NewCmd(container); err != nil {
		container.Logger.Fatal("failed to run commands", zap.Error(err))
	}

	container.MQTTDispatcher.ConsumeBlocking(container.Sig)
}
