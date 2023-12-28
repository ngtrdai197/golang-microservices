package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	HttpPort      int    `mapstructure:"HTTP_PORT" validate:"required"`
	GrpcPort      int    `mapstructure:"GRPC_PORT" validate:"required"`
	RedisHost     string `mapstructure:"REDIS_HOST" validate:"required"`
	RedisPort     int    `mapstructure:"REDIS_PORT" validate:"required"`
	RedisDatabase int    `mapstructure:"REDIS_DB" validate:"required"`
}

var Cfg *Config

func LoadConfig() {
	viper.SetConfigFile("config/.env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
	err = viper.Unmarshal(&Cfg)
	if err != nil {
		panic(err)
	}
	err = validator.New().Struct(Cfg)
	if err != nil {
		panic(err)
	}
}
