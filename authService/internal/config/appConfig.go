package config

import (
	"bytes"
	"gopkg.in/yaml.v2"
	"log"
	"os/exec"
	"strconv"
)

var JwtSecretKey []byte
var DBConfig DBSettings
var ServerPort string

type AppConfiguration struct {
	JWTSecretKey string         `yaml:"jwtSecretKey"`
	DB           DBSettings     `yaml:"db"`
	Server       ServerSettings `yaml:"server"`
}

type DBSettings struct {
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	DBName   string `yaml:"dbname"`
	SSLMode  string `yaml:"sslmode"`
}

type ServerSettings struct {
	Port int `yaml:"port"`
}

func InitAppConfig() {
	data, err := decryptConfigFile("config.yaml.gpg")
	if err != nil {
		log.Fatalf("Error decrypting config file: %v", err)
	}

	var config AppConfiguration
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatalf("Error parsing decrypted YAML file: %v", err)
	}

	JwtSecretKey = []byte(config.JWTSecretKey)
	DBConfig = config.DB
	ServerPort = strconv.Itoa(config.Server.Port)
}

func decryptConfigFile(filePath string) ([]byte, error) {
	cmd := exec.Command("gpg", "--decrypt", filePath)

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	return out.Bytes(), nil
}
