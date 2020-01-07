package server

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v2"
)

type args struct {
	ip      string
	port    string
	connstr string
}

func TestNewConfig(t *testing.T) {

	defaultConfig := Config{
		Ip:      DefaultIpAddr,
		Port:    DefaultPort,
		ConnStr: "postgres://127.0.0.1:docker@docker/test?sslmode=disable:",
	}

	tests := []struct {
		name       string
		args       args
		want       *Config
		wantErr    bool
		createFile bool
		clearFunc  func()
	}{
		{
			name:       "read from file",
			args:       newArgs(defaultConfig),
			want:       &defaultConfig,
			createFile: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.createFile {
				err := createConfigFile(&defaultConfig)
				require.NoError(t, err, "error creating config file")
			}

			got, err := NewConfig(tt.args.ip, tt.args.port, tt.args.connstr)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
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

	return ioutil.WriteFile(configFileName, b, 0666)
}

func newArgs(conf Config) args {
	return args{
		ip:      conf.Ip,
		port:    conf.Port,
		connstr: conf.ConnStr,
	}
}
