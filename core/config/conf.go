package config

import (
	"flag"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Server struct {
	HttpPort int `yaml:"http_port"`
}
type Other struct {
	DataBase string `yaml:"database"`
}
type Config struct {
	Server Server `yaml:"server"`
	Other  map[string]string
}

var (
	Conf     *Config
	confFile = flag.String("conf", "./conf/app.yaml", "conf file name")
	Viper    *viper.Viper
)

func (conf *Config) InitConf() (*Config, error) {
	flag.Parse()

	data, err := ioutil.ReadFile(*confFile)
	if err != nil {
		return nil, err
	}
	Viper = viper.New()
	err = yaml.Unmarshal([]byte(data), conf)
	if err != nil {
		return nil, err
	}

	for k, v := range conf.Other {
		Viper.Set(k, readYaml(v))
	}
	return conf, nil
}

func readYaml(filename string) (m interface{}) {
	data, _ := ioutil.ReadFile(filename)
	m = make(map[interface{}]interface{})
	err := yaml.Unmarshal([]byte(data), &m)
	if err != nil {
		return nil
	}
	return m
}
