package server

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/caledfwlch1/enlabtest/tools"
	"gopkg.in/yaml.v2"
)

func NewConfig(ip, port, host, user, pass, database, options string, save bool) (*Config, error) {
	conf, err := readConfigFromFile()
	if err != nil {
		return nil, err
	}

	conf.Ip = tools.IIF(ip != "", ip, conf.Ip).(string)
	conf.Port = tools.IIF(port != "", port, conf.Port).(string)
	conf.Host = tools.IIF(host != "", host, conf.Host).(string)
	conf.User = tools.IIF(user != "", user, conf.User).(string)
	conf.Pass = tools.IIF(pass != "", pass, conf.Pass).(string)
	conf.Database = tools.IIF(database != "", database, conf.Database).(string)
	conf.Options = tools.IIF(options != "", options, conf.Options).(string)

	if save {
		saveConfigToFile(conf)
	}

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
