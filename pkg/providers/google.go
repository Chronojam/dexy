package providers

import (
	"fmt"
	"log"
	"net/http"

	"google.golang.org/api/admin/directory/v1"
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
func (g *GoogleApps) Finalise(c *http.Client, m *FinaliseMetadata) error {
	srv, err := admin.New(c)
	if err != nil {
		log.Fatalf("Unable to retrieve directory Client %v", err)
	}

	r, err := srv.Groups.List().UserKey(m.Subject).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve groups in domain.", err)
	}
	fmt.Println(r.Groups)
	return nil
}

func (g *GoogleApps) Name() string {
	return "googleapps"
}
