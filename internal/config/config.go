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
	Name      string       `mapstructure:"name"`
	Chain     string       `mapstructure:"chain"`
	Provider  string       `mapstructure:"provider"`
	Endpoint  string       `mapstructure:"endpoint"`
	Protocol  string       `mapstructure:"protocol"`
	RateLimit []*RateLimit `mapstructure:"rate_limit"`
}

type RateLimit struct {
	TimeWindow time.Duration `mapstructure:"time_window"`
	Limit      uint64        `mapstructure:"limit"`
}

func ReadConfig(file string) (*Config, error) {
	v := viper.New()
	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var cfg Config
	if err := v.Unmarshal(&cfg); err != nil {
		return nil, err
	}
	cfg.LoadedFile = v.ConfigFileUsed()
	return &cfg, nil
}
