package api

import (
	"flag"
	"os"
	"time"
)

type delays struct {
	invoiceProcessing time.Duration
	orderProcessing   time.Duration
}

type Config struct {
	serverAddress        string
	databaseDSN          string
	accrualSystemAddress string
	jwtSecretKey         string
	delays               delays
}

func InitConfig() *Config {
	config := Config{
		delays: delays{
			invoiceProcessing: time.Duration(1),
			orderProcessing:   time.Duration(1),
		},
	}
	flag.StringVar(&config.serverAddress, "a", "0.0.0.0:8080", "server address")
	flag.StringVar(&config.databaseDSN, "d", "postgresql://loyalty:loyalty@192.168.33.26:5432/loyalty?sslmode=disable", "database URI")
	flag.StringVar(&config.accrualSystemAddress, "r", "", "accrual system address")
	flag.StringVar(&config.jwtSecretKey, "k", "very secret key", "secret key for JWT")
	flag.Parse()

	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		config.serverAddress = envRunAddress
	}

	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		config.databaseDSN = envDatabaseURI
	}

	if envAccrualSystemAddress := os.Getenv("ACCRULA_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		config.accrualSystemAddress = envAccrualSystemAddress
	}

	if envJWTSecretKey := os.Getenv("JWT_SECRET_KEY"); envJWTSecretKey != "" {
		config.jwtSecretKey = envJWTSecretKey
	}

	return &config
}
