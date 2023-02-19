package config

import (
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	LoadedFile string
	Webserver  *Webserver `mapstructure:"webserver"`
	Balancer   *Balancer  `mapstructure:"balancer"`
	Nodes      []*Node    `mapstructure:"nodes"`
}

type Webserver struct {
	Addr        string        `mapstructure:"addr"`
	ReadTimeout time.Duration `mapstructure:"read_timeout"`
}

type Balancer struct {
	TotalTimeout time.Duration `mapstructure:"total_timeout"`
	NodeTimeOut  time.Duration `mapstructure:"node_timeout"`
}

type Node struct {
	Name     string `mapstructure:"name"`
	Chain    string `mapstructure:"chain"`
	Provider string `mapstructure:"provider"`
	Endpoint string `mapstructure:"endpoint"`
	Protocol string `mapstructure:"protocol"`
}

func Read(file string) (*Config, error) {
	viper.SetConfigFile(file)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	cfg.LoadedFile = viper.ConfigFileUsed()
	return &cfg, nil
}
