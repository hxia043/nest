package main

import (
	"os"

	"github.com/hxia043/nest/internal/swicctl/cmd"
)

func main() {
	command := cmd.NewSWICCtlExecute()
	if err := command.Execute(); err != nil {
		os.Exit(1)
	}
}
