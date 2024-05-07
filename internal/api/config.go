package api

import (
	"flag"
	"os"
)

type Config struct {
	serverAddress        string
	databaseDSN          string
	accrualSystemAddress string
}

func InitConfig() *Config {
	config := Config{}
	flag.StringVar(&config.serverAddress, "a", "0.0.0.0:8080", "server address")
	flag.StringVar(&config.databaseDSN, "d", "", "database URI")
	flag.StringVar(&config.accrualSystemAddress, "r", "", "accrual system address")
	flag.Parse()

	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		config.serverAddress = envRunAddress
	}

	if envDatabaseUri := os.Getenv("DATABASE_URI"); envDatabaseUri != "" {
		config.databaseDSN = envDatabaseUri
	}

	if envAccrualSystemAddress := os.Getenv("ACCRULA_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		config.accrualSystemAddress = envAccrualSystemAddress
	}

	return &config
}
