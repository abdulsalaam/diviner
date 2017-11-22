package config

import "github.com/spf13/viper"

var config *viper.Viper

func init() {
	config = viper.New()
	config.AddConfigPath(".")
	config.SetConfigName("config")
}

// Load load a default configuration file
func Load() (*viper.Viper, error) {
	if err := config.ReadInConfig(); err != nil {
		return nil, err
	}

	return config, nil
}

// LoadFile create a new Viper to load specific configuration file
func LoadFile(config string) (*viper.Viper, error) {
	v := viper.New()

	v.SetConfigFile(config)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	return v, nil
}
