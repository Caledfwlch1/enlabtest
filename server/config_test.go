package server

import (
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"gopkg.in/yaml.v2"

	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	type args struct {
		ip   string
		port string
	}

	defaultConfig := &Config{
		Ip:   DefaultIpAddr,
		Port: DefaultPort,
	}
	wantConfig := &Config{
		Ip:   "127.0.0.1",
		Port: "1234",
	}

	tests := []struct {
		name       string
		args       args
		want       *Config
		createFile bool
		clearFunc  func()
	}{{
		name: "default config",
		want: defaultConfig,
	}, {
		name: "config from command line ",
		args: args{ip: "127.0.0.1", port: "1234"},
		want: wantConfig,
	}, {
		name:       "config from file",
		args:       args{ip: "127.0.0.1", port: "1234"},
		want:       wantConfig,
		createFile: true,
		clearFunc: func() {
			_ = os.Remove(configFileName)
		},
	}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createFile {
				err := createConfigFile(wantConfig)
				require.NoError(t, err, "error creating config file")
			}
			if got := NewConfig(tt.args.ip, tt.args.port); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewConfig() = %v, want %v", got, tt.want)
			}
			if tt.clearFunc != nil {
				tt.clearFunc()
			}
		})
	}
}

func createConfigFile(conf *Config) error {
	if conf == nil {
		return fmt.Errorf("configuration data is empty")
	}

	b, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(configFileName, b, 0x666)
}
