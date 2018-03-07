package providers

import (
	"log"
)

type Config struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	CallbackURL  string `mapstructure:"callback_url"`

	Provider interface{} `mapstructure:"provider"`
}

type IProvider interface {
	Finalise() error
	BuildRequestParameters() (string, error)
	IsSetCorrectly() (bool, error)
}

type BaseProvider struct {
}

func (provider BaseProvider) IsSetCorrectly() (bool, error) {
	return false, nil
}

func (provider BaseProvider) BuildRequestParameters() (string, error) {
	return "", nil
}

func (provider BaseProvider) Finalise() error {
	log.Println("A")

	return nil
}
