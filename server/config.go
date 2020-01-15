package server

import (
	"fmt"
	"io/ioutil"

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

func ifString(cond bool, outTrue, outFalse string) string {
	if cond {
		return outTrue
	}
	return outFalse
}
