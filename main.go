package main

import (
	"errors"
	"log"
	"os"

	fuse "github.com/datasance/flowfuse-agent/internal/fuse"
	sdk "github.com/datasance/iofog-go-sdk/v3/pkg/microservices"
)

var (
	flowFuseAgent *fuse.Agent
)

func init() {
	flowFuseAgent = new(fuse.Agent)
	flowFuseAgent.Config = new(fuse.Config)
}

func main() {
	ioFogClient, clientError := sdk.NewDefaultIoFogClient()
	if clientError != nil {
		log.Fatalln(clientError.Error())
	}

	// Update initial configuration
	if err := updateConfig(ioFogClient, flowFuseAgent.Config); err != nil {
		log.Fatalln(err.Error())
	}

	// Establish WebSocket connection for configuration updates
	confChannel := ioFogClient.EstablishControlWsConnection(0)

	// Channel for server exit handling
	exitChannel := make(chan error)

	// Start flowFuseAgent agent in a goroutine
	go flowFuseAgent.StartAgent(flowFuseAgent.Config, exitChannel)

	// Main loop to handle configuration updates
	for {
		select {
		case <-exitChannel:
			os.Exit(0)
		case <-confChannel:
			newConfig := new(fuse.Config)
			if err := updateConfig(ioFogClient, newConfig); err != nil {
				log.Fatal(err)
			} else {
				flowFuseAgent.UpdateAgent(newConfig)
			}
		}
	}
}

func updateConfig(ioFogClient *sdk.IoFogClient, config interface{}) error {
	attemptLimit := 5
	var err error

	for err = ioFogClient.GetConfigIntoStruct(config); err != nil && attemptLimit > 0; attemptLimit-- {
		return err
	}

	if attemptLimit == 0 {
		return errors.New("Update config failed")
	}

	return nil
}
