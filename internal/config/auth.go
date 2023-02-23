package config

import (
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
)

// Auth main auth config structure
type Auth struct {
	APIKeys []*APIKey `mapstructure:"api_keys"`

	//
	keys map[string]*APIKey
}

func (a *Auth) reIndexKeys() {
	a.keys = map[string]*APIKey{}
	if len(a.APIKeys) > 0 {
		for _, key := range a.APIKeys {
			key.RateLimiter = rate.NewLimiter(rate.Every(key.RateLimit.Per), key.RateLimit.Rate)
			a.keys[key.Key] = key
		}
	}
}

// Auth main authentication method, checks if key is valid & get api info
func (a *Auth) Auth(key string) (apikey *APIKey, found bool) {
	apiKey, found := a.keys[key]
	return apiKey, found
}

// APIKey main apikey structure (user)
type APIKey struct {
	Key          string    `mapstructure:"key"`
	RateLimit    RateLimit `mapstructure:"rate_limit"`
	AllowedHosts []string  `mapstructure:"allowed_hosts"`
	//todo: create rate limiter instances here
	RateLimiter *rate.Limiter
}

// ReadAuth read auth (user) info from a yaml file
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
