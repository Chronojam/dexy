package providers

import (
	"fmt"
	"log"
)

type GoogleApps struct {
	IProvider
	HD string
}

func (g GoogleApps) IsSetCorrectly() (bool, error) {
	log.Println("GAPPS isSet")
	return true, nil
}

func (g GoogleApps) BuildRequestParameters() (string, error) {
	log.Println("GAPPS brp")

	return fmt.Sprintf("&hd=%s", g.HD), nil
}
func (g GoogleApps) Finalise() error {
	log.Println("B")
	return nil
}
