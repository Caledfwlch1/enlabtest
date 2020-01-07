package server

import (
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/yaml.v2"
)

func NewConfig(ip, port, connStr string) (*Config, error) {
	conf, err := readConfigFromFile()
	if err != nil {
		return nil, err
	}

	conf.Ip = ifString(ip != "", ip, conf.Ip)
	conf.Port = ifString(port != "", port, conf.Port)
	conf.ConnStr = ifString(connStr != "", connStr, conf.ConnStr)

	return conf, nil
}

func readConfigFromFile() (*Config, error) {
	var conf Config
	b, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return nil, fmt.Errorf("error reading config file %s : %s", configFileName, err)
	}

	if yaml.Unmarshal(b, &conf) != nil {
		return nil, fmt.Errorf("error parsing config file %s : %s", configFileName, err)
	}

	return &conf, nil
}

func saveConfigToFile(conf *Config) {
	b, err := yaml.Marshal(conf)
	if err != nil {
		log.Println("configuration data encoding error")
		return
	}

	err = ioutil.WriteFile(configFileName, b, 0x666)
	if err != nil {
		log.Println("error writing to configuration file")
	}
}

func ifString(cond bool, outTrue, outFalse string) string {
	if cond {
		return outTrue
	}
	return outFalse
}
