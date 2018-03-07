package providers

import (
	"fmt"
	"log"
)

type GoogleApps struct {
	HD     string
	Scopes []string
}

func (g *GoogleApps) Validate(model map[string]interface{}) (bool, error) {
	g.HD = model["hd"].(string)
	if s, ok := model["scopes"]; ok {
		scopes := s.([]interface{})
		g.Scopes = make([]string, len(scopes))
		for i, v := range scopes {
			g.Scopes[i] = fmt.Sprintf(v.(string))
		}
	}
	return true, nil
}

func (g *GoogleApps) AdditionalScopes() []string {
	return g.Scopes
}

func (g *GoogleApps) BuildRequestParameters() (string, error) {
	return fmt.Sprintf("&hd=%s", g.HD), nil
}
func (g *GoogleApps) Finalise() error {
	log.Println("B")
	return nil
}

func (g *GoogleApps) Name() string {
	return "googleapps"
}
