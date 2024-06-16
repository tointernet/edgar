package cmd

import "github.com/tointernet/edgar/pkgs"

func NewCmd(container *pkgs.Container) error {
	if err := NewMQTTConsumer(container); err != nil {
		return err
	}

	return nil
}
