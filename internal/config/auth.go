package config

import "github.com/spf13/viper"

type Auth struct {
	APIKeys []*APIKey `mapstructure:"api_keys"`

	//
	keys map[string]*APIKey
}

func (a *Auth) reIndexKeys() {
	a.keys = map[string]*APIKey{}
	if len(a.APIKeys) > 0 {
		for _, key := range a.APIKeys {
			a.keys[key.Key] = key
		}
	}
}

func (a *Auth) Auth(key string) (apikey *APIKey, found bool) {
	apiKey, found := a.keys[key]
	return apiKey, found
}

type APIKey struct {
	Key          string      `mapstructure:"key"`
	RateLimit    []RateLimit `mapstructure:"rate_limit"`
	AllowedHosts []string    `mapstructure:"allowed_hosts"`
	//todo: create rate limiter instances here
}

func ReadAuth(file string) (*Auth, error) {
	v := viper.New()
	v.SetConfigFile(file)
	if err := v.ReadInConfig(); err != nil {
		return nil, err
	}
	var auth Auth
	if err := v.Unmarshal(&auth); err != nil {
		return nil, err
	}
	auth.reIndexKeys()
	return &auth, nil
}
