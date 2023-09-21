package config

import "github.com/spf13/viper"

type Config struct {
	PostgresDriver string `mapstructure:"POSTGRES_DRIVER"`
	PostgresConn   string `mapstructure:"POSTGRES_CONNECTION"`

	Port string `mapstructure:"PORT"`
}

func LoadConfig(path string) (config Config, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigName("app")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
