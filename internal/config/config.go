package config

import (
	"flag"
	"log/slog"
	"os"
)

type MainConfig struct {
	RunAddress  string
	DatabaseURI string
}

func MakeConfig() MainConfig {
	config := MainConfig{
		RunAddress:  "0.0.0.0:8080",
		DatabaseURI: "",
	}

	return config
}

func (c *MainConfig) InitConfig() {
	// c.InitFlags()
	c.Parse()
}

func (c *MainConfig) InitFlags() {
	flag.StringVar(&c.RunAddress, "a", "localhost:8080", "default host for server")
	flag.StringVar(&c.DatabaseURI, "d", "", "database URI")

	slog.Info("flags inited")
}

func (c *MainConfig) Parse() {
	flag.Parse()

	if e := os.Getenv("RUN_ADDRESS"); e != "" {
		c.RunAddress = e
	}
	if e := os.Getenv("DATABASE_URI"); e != "" {
		c.DatabaseURI = e
	}
}
