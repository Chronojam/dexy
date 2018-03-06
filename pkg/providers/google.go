package providers

import "log"

type GoogleApps struct {
	BaseProvider
	HD string
}

func (g GoogleApps) IsSetCorrectly() (bool, error) {
	log.Println("GAPPS isSet")
	return true, nil
}

func (g GoogleApps) Finalise() error {
	log.Println("B")
	return nil
}
