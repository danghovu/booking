package config

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Server struct {
		Port string `mapstructure:"port"`
	} `mapstructure:"server"`
	Postgres struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		User     string `mapstructure:"user"`
		Password string `mapstructure:"password"`
		DbName   string `mapstructure:"db_name"`
		SSLMode  string `mapstructure:"ssl_mode"`
		OpenConn int    `mapstructure:"open_conn"`
		IdleConn int    `mapstructure:"idle_conn"`
	} `mapstructure:"postgres"`
	Redis struct {
		Host     string `mapstructure:"host"`
		Port     int    `mapstructure:"port"`
		Password string `mapstructure:"password"`
		Prefix   string `mapstructure:"prefix"`
	} `mapstructure:"redis"`
	Booking struct {
		MaxBookingPerUser int `mapstructure:"max_booking_per_user"`
	} `mapstructure:"booking"`
	Token struct {
		LockedDuration time.Duration `mapstructure:"locked_duration"`
	} `mapstructure:"token"`
	JWT struct {
		SecretKey       string        `mapstructure:"secret_key"`
		AccessTokenExp  time.Duration `mapstructure:"access_token_exp"`
		RefreshTokenExp time.Duration `mapstructure:"refresh_token_exp"`
	} `mapstructure:"jwt"`
	SupportingMoney struct {
		Currency string `mapstructure:"currency"`
	} `mapstructure:"supporting_money"`
	GracefulShutdown time.Duration `mapstructure:"graceful_shutdown"`
	Asynq            struct {
		Concurrency int            `mapstructure:"concurrency"`
		Queues      map[string]int `mapstructure:"queues"`
	} `mapstructure:"asynq"`
}

func NewConfig() (*Config, error) {
	v := viper.New()

	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("./config")

	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	v.AutomaticEnv()
	var c Config
	if err := v.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unable to decode into struct: %w", err)
	}

	log.Printf("config: %+v", c)

	return &c, nil
}
