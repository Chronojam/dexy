package providers

import "net/http"

type Config struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	CallbackURL  string `mapstructure:"callback_url"`
}

type FinaliseMetadata struct {
	Subject string
}

type IProvider interface {
	Finalise(*http.Client, *FinaliseMetadata) error
	BuildRequestParameters() (string, error)
	Validate(map[string]interface{}) (bool, error)
	AdditionalScopes() []string
	Name() string
}
