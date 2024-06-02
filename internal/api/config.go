package api

import (
	"flag"
	"os"
	"time"
)

type delays struct {
	ordersCalculation time.Duration
	ordersProcessing  time.Duration
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
			ordersCalculation: time.Duration(500) * time.Millisecond,
			ordersProcessing:  time.Duration(500) * time.Millisecond,
		},
	}
	flag.StringVar(&config.serverAddress, "a", "0.0.0.0:80", "server address")
	flag.StringVar(&config.databaseDSN, "d", "", "database URI")
	flag.StringVar(&config.accrualSystemAddress, "r", "http://127.0.0.1:8080", "accrual system address")
	flag.StringVar(&config.jwtSecretKey, "k", "very secret key", "secret key for JWT")
	flag.Parse()

	if envRunAddress := os.Getenv("RUN_ADDRESS"); envRunAddress != "" {
		config.serverAddress = envRunAddress
	}

	if envDatabaseURI := os.Getenv("DATABASE_URI"); envDatabaseURI != "" {
		config.databaseDSN = envDatabaseURI
	}

	if envAccrualSystemAddress := os.Getenv("ACCRUAL_SYSTEM_ADDRESS"); envAccrualSystemAddress != "" {
		config.accrualSystemAddress = envAccrualSystemAddress
	}

	if envJWTSecretKey := os.Getenv("JWT_SECRET_KEY"); envJWTSecretKey != "" {
		config.jwtSecretKey = envJWTSecretKey
	}

	return &config
}
