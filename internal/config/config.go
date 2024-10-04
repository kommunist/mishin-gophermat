package config

import (
	"flag"
	"log/slog"
	"os"
)

type MainConfig struct {
	RunAddress  string
	DatabaseURI string
	AccrualURI  string
}

func MakeConfig() MainConfig {
	config := MainConfig{
		RunAddress:  "",
		DatabaseURI: "",
		AccrualURI:  "",
	}

	return config
}

func (c *MainConfig) InitConfig() {
	c.ParseEnv()
	c.InitFlags()
}

func (c *MainConfig) InitFlags() {
	var f string
	flag.StringVar(&f, "a", "localhost:8080", "default host for server")
	if c.RunAddress == "" {
		c.RunAddress = f
	}

	flag.StringVar(&f, "d", "", "database URI")
	if c.DatabaseURI == "" {
		c.DatabaseURI = f
	}

	flag.StringVar(&f, "r", "", "accrual URI")
	if c.AccrualURI == "" {
		c.AccrualURI = f
	}

	slog.Info("flags inited")
}

func (c *MainConfig) ParseEnv() {
	flag.Parse()

	if e := os.Getenv("RUN_ADDRESS"); e != "" {
		c.RunAddress = e
	}
	if e := os.Getenv("DATABASE_URI"); e != "" {
		c.DatabaseURI = e
	}
	if e := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); e != "" {
		c.AccrualURI = e
	}
	slog.Info("env parsed")
}
