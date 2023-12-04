package config

import (
	"github.com/joho/godotenv"
	"os"
	"strconv"
)

type (
	Config struct {
		Postgres    Postgres    `json:"postgres"`
		HTTTPServer HTTTPServer `json:"http_server"`
	}

	Postgres struct {
		URL string `json:"url"`
	}

	HTTTPServer struct {
		Hostname   string `json:"hostname"`
		Port       string `json:"port"`
		TypeServer string `json:"type_server"`
	}
)

func New() (*Config, error) {
	err := godotenv.Load("configs/server.env")
	if err != nil {
		return nil, err
	}

	config := &Config{
		Postgres: Postgres{
			URL: os.Getenv("POSTGRES_URL"),
		},
		HTTTPServer: HTTTPServer{
			Hostname:   os.Getenv("HTTP_SERVER_HOSTNAME"),
			Port:       os.Getenv("HTTP_SERVER_PORT"),
			TypeServer: os.Getenv("HTTP_SERVER_TYPE_SERVER"),
		},
	}

	return config, nil
}

func parseEnvInt(value string) int {
	intValue, err := strconv.Atoi(value)
	if err != nil {
		return 0
	}
	return intValue
}
