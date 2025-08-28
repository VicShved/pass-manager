package config

import (
	"encoding/json"
	"flag"
	"os"
)

// AuthorizationTokenName
const AuthorizationTokenName string = "Authorization"

// type ServerConfigStruct
type ServerConfigStruct struct {
	ServerAddress   string `json:"server_address"`
	ServerPort      string `json:"port"`
	FileStoragePath string `json:"file_storage_path"`
	DBDSN           string `json:"database_dsn"`
	SecretKey       string
	LogLevel        string
	EnableTLS       bool `json:"enable_tls"`
	ConfigFileName  string
	SchemaName      string
}

// var ServerConfig
var ServerConfig ServerConfigStruct

func getConfigArgsEnvVars() *ServerConfigStruct {
	flag.StringVar(&ServerConfig.ServerAddress, "a", "localhost", "start base url")
	flag.StringVar(&ServerConfig.ServerPort, "p", "7777", "port")
	flag.StringVar(&ServerConfig.FileStoragePath, "f", "", "file storage path")
	flag.BoolVar(&ServerConfig.EnableTLS, "s", false, "enable tls")
	flag.StringVar(&ServerConfig.DBDSN, "d", "", "DataBase DSN")
	flag.StringVar(&ServerConfig.SecretKey, "k", "VeryImpotantSecretKey.YesYes", "Secret key")
	flag.StringVar(&ServerConfig.LogLevel, "l", "INFO", "Log level")
	flag.StringVar(&ServerConfig.ConfigFileName, "c", "", "Config from file")
	flag.Parse()

	value, exists := os.LookupEnv("SERVER_ADDRESS")
	if exists {
		ServerConfig.ServerAddress = value
	}

	value, exists = os.LookupEnv("SERVER_PORT")
	if exists {
		ServerConfig.ServerPort = value
	}

	value, exists = os.LookupEnv("FILE_STORAGE_PATH")
	if exists {
		ServerConfig.FileStoragePath = value
	}

	value, exists = os.LookupEnv("DATABASE_DSN")
	if exists {
		ServerConfig.DBDSN = value
	}

	value, exists = os.LookupEnv("SECRET_KEY")
	if exists {
		ServerConfig.SecretKey = value
	}

	value, exists = os.LookupEnv("LOG_LEVEL")
	if exists {
		ServerConfig.LogLevel = value
	}
	return &ServerConfig
}

func readConfigFromFile(fileName string) (*ServerConfigStruct, error) {

	b, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	var config ServerConfigStruct
	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

func updateConfig(target *ServerConfigStruct, source *ServerConfigStruct) *ServerConfigStruct {
	if target.ServerAddress == "" {
		target.ServerAddress = source.ServerAddress
	}
	if target.ServerPort == "" {
		target.ServerPort = source.ServerPort
	}
	if target.FileStoragePath == "" {
		target.FileStoragePath = source.FileStoragePath
	}
	if target.DBDSN == "" {
		target.DBDSN = source.DBDSN
	}
	return target
}

// GetServerConfig return app config
func GetServerConfig() *ServerConfigStruct {
	serverConfig := getConfigArgsEnvVars()
	if serverConfig.ConfigFileName != "" {
		fileConfig, err := readConfigFromFile(serverConfig.ConfigFileName)
		if err == nil {
			updateConfig(&ServerConfig, fileConfig)
		}
	}
	return serverConfig
}
