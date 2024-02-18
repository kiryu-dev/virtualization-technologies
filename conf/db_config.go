package conf

import (
	"fmt"
	"os"
)

type DbConfig struct {
	Host     string
	Port     string
	Username string
	DBName   string
	Password string
	SslMode  string
}

func NewDbConfig() DbConfig {
	return DbConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Username: os.Getenv("DB_USERNAME"),
		DBName:   os.Getenv("DB_NAME"),
		Password: os.Getenv("DB_PASSWORD"),
		SslMode:  os.Getenv("DB_SSLMODE"),
	}
}

func (d DbConfig) String() string {
	return fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=%s",
		d.Username, d.Password, d.Host, d.Port, d.DBName, d.SslMode)
}
