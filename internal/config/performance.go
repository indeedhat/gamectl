package config

import (
	"fmt"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

type Performance struct {
	Mount   []string
	Network []string
}

// LoadPerformanceConfig Load the LoadPerformanceConfig from file
func LoadPerformanceConfig() (*Performance, error) {
	performance := Performance{}

	data, err := ioutil.ReadFile(fmt.Sprintf("%s/performance.yml", ConfigDirectory))
	if err != nil {
		return nil, err
	}

	if err = yaml.Unmarshal(data, &performance); err != nil {
		return nil, err
	}

	return &performance, nil
}
