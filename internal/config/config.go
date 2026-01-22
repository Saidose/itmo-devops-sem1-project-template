package config

import (
	"fmt"
	"os"
	"strings"
)

// Config содержит все конфигурационные параметры приложения
type Config struct {
	HttpPort string // Порт для HTTP сервера
	DBUser   string // URL для подключения к базе данных
	DBPass   string
	DBName   string
	DBHost   string
}

// mustGetEnv получает значение обязательной переменной окружения или возвращает ошибку если она пустая
func mustGetEnv(key string) (string, error) {
	value := os.Getenv(key)
	if value == "" {
		return "", fmt.Errorf("%s is required", key)
	}
	return value, nil
}

// LoadConfig загружает конфигурацию из переменных окружения и возвращает Config
// Возвращает ошибку если какие-то обязательные переменные не установлены
func LoadConfig() (*Config, error) {
	errs := make([]string, 0)

	httpPort, err := mustGetEnv("HTTP_PORT")
	if err != nil {
		errs = append(errs, err.Error())
	}

	dbUser, err := mustGetEnv("DB_USER")
	if err != nil {
		errs = append(errs, err.Error())
	}
	dbPass, err := mustGetEnv("DB_PASS")
	if err != nil {
		errs = append(errs, err.Error())
	}
	dbName, err := mustGetEnv("DB_NAME")
	if err != nil {
		errs = append(errs, err.Error())
	}
	dbHost, err := mustGetEnv("DB_HOST")
	if err != nil {
		errs = append(errs, err.Error())
	}

	if len(errs) > 0 {
		return nil, fmt.Errorf("config validation failed:\n  %s", strings.Join(errs, "\n  "))
	}
	return &Config{
		HttpPort: httpPort,
		DBUser:   dbUser,
		DBPass:   dbPass,
		DBName:   dbName,
		DBHost:   dbHost,
	}, nil
}
