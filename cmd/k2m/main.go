package main

import (
	"github.com/tointernet/edgar/pkg"
)

func main() {
	container, err := pkg.NewContainer()
	if err != nil {
		panic(err)
	}

	container.Logger.Debug("do something else...")

	err = container.TinyHTTPServer.Run()
	if err != nil {
		panic(err)
	}

	<-container.Sig
}
