package fuse

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"github.com/datasance/flowfuse-agent/internal/exec"
)

type Agent struct {
	Config *Config
	mu     sync.Mutex // Mutex to ensure that only one agent is started at a time
}

type Config struct {
	DeviceId         string `json:"deviceId"`
	Token            string `json:"token"`
	CredentialSecret string `json:"credentialSecret"`
	ForgeURL         string `json:"forgeURL"`
	BrokerURL        string `json:"brokerURL"`
	BrokerUsername   string `json:"brokerUsername"`
	BrokerPassword   string `json:"brokerPassword"`
}

func (a *Agent) UpdateAgent(config *Config) error {
	a.mu.Lock() // Ensure only one agent is started at a time
	defer a.mu.Unlock()

	// Create the new configuration files
	if err := a.createConfigFiles(config); err != nil {
		return err
	}

	log.Printf("FlowFuse agent configuration updated successfully.")
	return nil
}

func (a *Agent) createConfigFiles(config *Config) error {
	configDir := "./opt/flowfuse-device"
	log.Printf("Creating directory: %s", configDir)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	deviceConfPath := filepath.Join(configDir, "device.yml")
	log.Printf("Creating device config file at %s", deviceConfPath)
	if err := flowFuseDeviceFile(deviceConfPath, config); err != nil {
		return fmt.Errorf("failed to create account config file: %v", err)
	}

	log.Printf("FlowFuse Agent configuration files updated successfully in %s", configDir)
	return nil
}

func flowFuseDeviceFile(path string, config *Config) error {
	// Format the configuration in YAML-like structure
	data := fmt.Sprintf(`deviceId: %s
token: %s
credentialSecret: %s
forgeURL: %s
brokerURL: %s
brokerUsername: %s
brokerPassword: %s
`,
		config.DeviceId,
		config.Token,
		config.CredentialSecret,
		config.ForgeURL,
		config.BrokerURL,
		config.BrokerUsername,
		config.BrokerPassword,
	)

	// Write the configuration to the specified file path
	if err := ioutil.WriteFile(path, []byte(data), 0644); err != nil {
		return fmt.Errorf("failed to write device config file: %v", err)
	}

	log.Printf("Device configuration file written successfully at %s", path)
	return nil
}

func (a *Agent) StartAgent(config *Config, exitChannel chan error) error {

	// First, update the FlowFuse agent configuration and stop any existing agent if needed
	if err := a.UpdateAgent(config); err != nil {
		return fmt.Errorf("failed to update FlowFuse agent before starting: %v", err)
	}

	args := []string{}

	env := []string{} // Pass any required environment variables here

	go exec.Run(exitChannel, "flowfuse-device-agent", args, env)

	log.Printf("Flowfuse Agent started successfully")

	return nil
}
