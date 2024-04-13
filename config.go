package main

import "github.com/kelseyhightower/envconfig"

type Config struct {
	// General configs
	TemperatureUnits string `default:"Fahrenheit"`
	PollingInterval  int    `default:"5"`
	// Screen configs
	ScreenWidth         int  `default:"128"`
	ScreenHeight        int  `default:"32"`
	ScreenRotated       bool `default:"true"`
	ScreenSquential     bool `default:"true"`
	ScreenSwapTopBottom bool `default:"false"`
	ScreenActive        bool `default:"true"`
	// Prometheus Metrics
	PromActive bool `default:"true"`
}

func InitConfig(config *Config) {
	err := envconfig.Process("tempsense", &config)
	if err != nil {
		panic(err)
	}
}
