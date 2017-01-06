package configutils

import (
	"encoding/json"
	"fmt"
	"os"
)

// Loads configuration. If name of configuration file provided then reads from it otherwise
// tries to read from the specified environment variable
func LoadConfig(configFileName, configEnvVarName string, config interface{}) error {
	_ = "breakpoint"
	if len(configFileName) == 0 {
		return loadConfigFromEnv(configEnvVarName, config)
	}
	return loadConfigFromFile(configFileName, &config)
}

// LoadConfigFromFile loads configuration from the file with the specified name.
func loadConfigFromFile(configFileName string, config interface{}) error {
	if len(configFileName) == 0 {
		return fmt.Errorf("Name of configuration file cannot be an empty string.")
	}
	file, err := os.Open(configFileName)
	if err != nil {
		return err
	}
	return json.NewDecoder(file).Decode(&config)
}

func loadConfigFromEnv(envVarName string, config interface{}) error {
	if len(envVarName) == 0 {
		return fmt.Errorf("Name of env variable with config cannot be an empty string")
	}
	envVarValue := os.Getenv(envVarName)
	byt := []byte(envVarValue)
	return json.Unmarshal(byt, &config)
}
