package config

import "github.com/spf13/viper"

type Auth struct {
	APIKeys []APIKeys `mapstructure:"api_keys"`
}

type APIKeys struct {
	Key          string      `mapstructure:"key"`
	RateLimit    []RateLimit `mapstructure:"rate_limit"`
	AllowedHosts []string    `mapstructure:"allowed_hosts"`
}

func ReadAuth(file string) (*Auth, error) {
	v := viper.New()
	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Auth
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
