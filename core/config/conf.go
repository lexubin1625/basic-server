package config

import (
	"flag"
	"fmt"
	"github.com/spf13/viper"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Env uint8

const (
	Local int8 = iota
	Test
	Gray
	Pro
)

type Server struct {
	HttpPort      int  `yaml:"http_port"`
	SwaggerEnable bool `yaml:"swagger_enable"`
}
type Other struct {
	DataBase string `yaml:"database"`
}
type Config struct {
	Env    int8   `yaml:"env"`
	Server Server `yaml:"server"`
	Other  map[string]string
}

var (
	Conf     *Config
	confFile = flag.String("conf", "./conf/config_local.yaml", "conf file name")
	Viper    *viper.Viper
)

func (conf *Config) InitConf() (*Config, error) {
	flag.Parse()

	data, err := ioutil.ReadFile(*confFile)
	if err != nil {
		return nil, err
	}
	Viper = viper.New()
	err = yaml.Unmarshal(data, conf)
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
		panic(fmt.Sprintf("file %s ,err :%v", filename, err))
	}
	return m
}
