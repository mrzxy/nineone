package config

import (
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

var Conf *Config

type Config struct {
	Server struct {
		Port           string  `yaml:"Port"`     // 端口
		MySqlUrl       string  `yaml:"MySqlUrl"` // 数据库连接地址
		Debug          bool    `yaml:"Debug"`
	} `yaml:"Server"`

}

func InitConfig(filename string) {
	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		logrus.Error(err)
		return
	}

	Conf = &Config{}
	err = yaml.Unmarshal(yamlFile, Conf)
	if err != nil {
		logrus.Error(err)
	}
}
