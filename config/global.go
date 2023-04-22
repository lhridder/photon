package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Globalconfig struct {
	Debug bool
}

func LoadGlobal() (*Globalconfig, error) {
	var cfg *Globalconfig
	file, err := ioutil.ReadFile("config.yml")
	if err != nil {
		return nil, fmt.Errorf("failed to open config.yml: %s", err)
	}

	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config.yml: %s", err)
	}

	return cfg, nil
}
