package config

import (
	"encoding/json"
	"flag"
	"os"
)

// go build -ldflags "-X config.Version=v0.1 -X 'main.BuildTime=$(date +'%Y/%m/%d %H:%M:%S')'"

// AuthorizationTokenName
const AuthorizationTokenName string = "Authorization"

// ClientConfigStruct
type ClientConfigStruct struct {
	ServerAddress  string `json:"server_address"`
	ConfigFileName string
	LogLevel       string `json:"log_level"`
}

// ClientConfig
var ClientConfig ClientConfigStruct

func getConfigArgsEnvVars() *ClientConfigStruct {
	flag.StringVar(&ClientConfig.ServerAddress, "a", "localhost:7777", "server url:port")
	flag.StringVar(&ClientConfig.ConfigFileName, "f", "", "config from file")
	flag.StringVar(&ClientConfig.LogLevel, "l", "INFO", "Log level")
	flag.Parse()

	value, exists := os.LookupEnv("SERVER_ADDRESS")
	if exists {
		ClientConfig.ServerAddress = value
	}
	return &ClientConfig
}

func readConfigFromFile(fileName string) (*ClientConfigStruct, error) {

	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var config ClientConfigStruct
	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func updateConfig(target *ClientConfigStruct, source *ClientConfigStruct) *ClientConfigStruct {
	if target.ServerAddress == "" {
		target.ServerAddress = source.ServerAddress
	}
	return target
}

// GetServerConfig
func GetClientConfig() *ClientConfigStruct {
	clientConfig := getConfigArgsEnvVars()
	if clientConfig.ConfigFileName != "" {
		fileConfig, err := readConfigFromFile(clientConfig.ConfigFileName)
		if err == nil {
			updateConfig(&ClientConfig, fileConfig)
		}
	}
	return clientConfig
}
