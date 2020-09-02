package config

import (
	"flag"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

type Server struct {
	HttpPort int `yaml:"http_port"`
}
type Config struct {
	Server Server `yaml:"server"`
}
var (
	Conf *Config
	confFile = flag.String("conf","./conf/app.yaml","conf file name")
)

func (conf *Config)InitConf()(*Config, error){
	flag.Parse()

	data, err := ioutil.ReadFile(*confFile)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal([]byte(data),conf)
	if err != nil {
		return nil, err
	}
	return Conf,nil
}