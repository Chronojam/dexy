package providers

type Config struct {
	ClientID     string `mapstructure:"client_id"`
	ClientSecret string `mapstructure:"client_secret"`
	CallbackURL  string `mapstructure:"callback_url"`
}

type IProvider interface {
	Finalise() error
	BuildRequestParameters() (string, error)
	Validate(map[string]interface{}) (bool, error)
	AdditionalScopes() []string
	Name() string
}
