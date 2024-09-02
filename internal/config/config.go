package config

import (
	"flag"
	"log/slog"
	"os"
)

type MainConfig struct {
	BaseServerURL string
	DatabaseDSN   string
}

func MakeConfig() MainConfig {
	config := MainConfig{
		BaseServerURL: "0.0.0.0:8080",
		DatabaseDSN:   "",
	}

	return config
}

func (c *MainConfig) InitConfig() {
	c.InitFlags()
	c.Parse()
}

func (c *MainConfig) InitFlags() {
	flag.StringVar(&c.BaseServerURL, "a", "localhost:8080", "default host for server")
	flag.StringVar(&c.DatabaseDSN, "d", "", "database DSN")

	slog.Info("flags inited")
}

func (c *MainConfig) Parse() {
	flag.Parse()

	if e := os.Getenv("SERVER_ADDRESS"); e != "" {
		c.BaseServerURL = e
	}
	if e := os.Getenv("DATABASE_DSN"); e != "" {
		c.DatabaseDSN = e
	}
}
