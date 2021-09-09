package common

import (
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"gopkg.in/yaml.v2"
)

type VirtualHost struct {
	Name    string `yaml:"name"`
	Path    string `yaml:"path"`
	Default string `yaml:"default"`
}

type VirtualHosts []VirtualHost

type Config struct {
	IP      string        `yaml:"ip"`
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
	Default string        `yaml:"default"`
	VHosts  VirtualHosts  `yaml:"virtualhosts"`
}

func (c *Config) Parse(data []byte) error {
	return yaml.Unmarshal(data, c)
}

func ReadConfig(configFile string) (Config, error) {
	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		panic(fmt.Sprintf("Couldn't read the config: %v\n", err))
	}

	var config Config
	err = config.Parse(data)
	if err != nil {
		return Config{}, err
	}
	log.Printf("Struct: %v\n", config)

	return config, nil
}
