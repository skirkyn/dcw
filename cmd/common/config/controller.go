package config

import "github.com/unknownfeature/dcw/cmd/util"

type ControllerConfig struct {
	Workers              int `json:"workers"`
	MaxSendRetries       int `json:"maxSendRetries"`
	MaxSendRetriesTtsSec int `json:"maxSendRetriesTtsSec"`
	ResLength            int `json:"resLength"`
}

func ReadControllerConfig() (ControllerConfig, error) {

	return util.ReadToStruct[ControllerConfig](configNames[Controller], func() ControllerConfig { return ControllerConfig{} })

}
