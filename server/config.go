package server

import (
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func NewConfig(ip, port string) *Config {
	var conf *Config

	printState := func(conf *Config) {
		log.Printf("using config: %s", conf)
	}

	if ip != "" && port != "" {
		conf = &Config{
			Ip:   ip,
			Port: port,
		}
		printState(conf)
		return conf
	}

	defaultConfig := &Config{
		Ip:   DefaultIpAddr,
		Port: DefaultPort,
	}

	b, err := ioutil.ReadFile(configFileName)
	if err != nil {
		conf = defaultConfig
		log.Printf("error reading config file %s : %s", configFileName, err)
		printState(conf)
		return conf
	}

	if yaml.Unmarshal(b, conf) != nil {
		conf = defaultConfig
		log.Printf("error parsing config file %s : %s", configFileName, err)
		printState(conf)
		return conf
	}

	return conf
}
