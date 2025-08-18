package config

import (
	"encoding/json"
	"flag"
	"os"
)

// type ServerConfigStruct
type ServerConfigStruct struct {
	ServerAddress   string `json:"server_address"`
	BaseURL         string `json:"base_url"`
	FileStoragePath string `json:"file_storage_path"`
	DBDSN           string `json:"database_dsn"`
	SecretKey       string
	LogLevel        string
	EnableTLS       bool `json:"enable_tls"`
	ConfigFileName  string
	// TrustedSubnet   string `json:"trusted_subnet"`
}

// var ServerConfig
var ServerConfig ServerConfigStruct

func getConfigArgsEnvVars() *ServerConfigStruct {
	flag.StringVar(&ServerConfig.ServerAddress, "a", "localhost:8080", "start base url")
	// flag.StringVar(&ServerConfig.BaseURL, "b", "http://localhost:8080", "result base url")
	flag.StringVar(&ServerConfig.FileStoragePath, "f", "", "file storage path")
	flag.BoolVar(&ServerConfig.EnableTLS, "s", false, "enable tls")
	flag.StringVar(&ServerConfig.DBDSN, "d", "", "DataBase DSN")
	flag.StringVar(&ServerConfig.SecretKey, "k", "VeryImpotantSecretKey.YesYes", "Secret key")
	flag.StringVar(&ServerConfig.LogLevel, "l", "INFO", "Log level")
	// flag.StringVar(&ServerConfig.TrustedSubnet, "t", "", "Trusted subnet")
	flag.Parse()

	value, exists := os.LookupEnv("SERVER_ADDRESS")
	if exists {
		ServerConfig.ServerAddress = value
	}

	value, exists = os.LookupEnv("BASE_URL")
	if exists {
		ServerConfig.BaseURL = value
	}

	value, exists = os.LookupEnv("FILE_STORAGE_PATH")
	if exists {
		ServerConfig.FileStoragePath = value
	}

	// value, exists = os.LookupEnv("DATABASE_DSN")
	// if exists {
	// 	ServerConfig.DBDSN = value
	// }

	value, exists = os.LookupEnv("SECRET_KEY")
	if exists {
		ServerConfig.SecretKey = value
	}

	value, exists = os.LookupEnv("LOG_LEVEL")
	if exists {
		ServerConfig.LogLevel = value
	}

	// value, exists = os.LookupEnv("TRUSTED_SUBNET")
	// if exists {
	// 	ServerConfig.TrustedSubnet = value
	// }

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
	if target.BaseURL == "" {
		target.BaseURL = source.BaseURL
	}
	if target.FileStoragePath == "" {
		target.FileStoragePath = source.FileStoragePath
	}
	if target.DBDSN == "" {
		target.DBDSN = source.DBDSN
	}
	return target
}

// func GetServerConfig
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
