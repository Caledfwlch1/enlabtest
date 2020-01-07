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
	ip       string
	port     string
	host     string
	user     string
	pass     string
	database string
	options  string
	save     bool
}

func TestNewConfig(t *testing.T) {

	defaultConfig := Config{
		Ip:       DefaultIpAddr,
		Port:     DefaultPort,
		Host:     "127.0.0.1",
		User:     "docker",
		Pass:     "docker",
		Database: "test",
		Options:  "sslmode=disable",
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
			args:       newArgs(defaultConfig, true),
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

			got, err := NewConfig(tt.args.ip, tt.args.port, tt.args.host, tt.args.user, tt.args.pass, tt.args.database, tt.args.options, tt.args.save)
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

func newArgs(conf Config, save bool) args {
	return args{
		ip:       conf.Ip,
		port:     conf.Port,
		host:     conf.Host,
		user:     conf.User,
		pass:     conf.Pass,
		database: conf.Database,
		options:  conf.Options,
		save:     save,
	}
}
