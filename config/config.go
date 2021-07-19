package config

import "github.com/spf13/viper"

type Config struct {
	DBFilePath string `mapstructure:"DB_FILEPATH"`
}

// Read configuration file and map it to a Config struct
func SetConfig() Config {
	viper.New()
	viper.SetConfigFile("./app.env")

	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}

	var configuration Config

	if err := viper.Unmarshal(&configuration); err != nil {
		panic(err)
	}

	return configuration
}
