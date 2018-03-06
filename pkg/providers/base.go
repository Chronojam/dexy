package providers

import (
	"log"
)

type Config struct {
	BaseProvider
	GoogleApps
}

type IProvider interface {
	Finalise() error
	IsSetCorrectly() (bool, error)

	GetClientID() string
	GetClientSecret() string
	GetCallbackURL() string
}

func (cmd IProvider) UnmarshalYAML(unmarshal func(interface{}) error) error {
	if err := unmarshal(&cmd); err != nil {
		return err
	}
	// Simply overwrite unmarshalled data PoC
	*cmd = []string{"overwritten", "values"}
	return nil
}

type BaseProvider struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	CallbackURL  string `mapstructure:"callback_url"`
}

func (provider BaseProvider) GetCallbackURL() string {
	return provider.CallbackURL
}
func (provider BaseProvider) GetClientID() string     { return provider.ClientID }
func (provider BaseProvider) GetClientSecret() string { return provider.ClientSecret }

func (provider BaseProvider) IsSetCorrectly() (bool, error) {
	return false, nil
}

func (provider BaseProvider) Finalise() error {
	log.Println("A")

	return nil
}
