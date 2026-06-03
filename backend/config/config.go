package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	QiniuKodo QiniuKodoConfig `yaml:"qiniu-kodo"`
	MySQL     MySQLConfig     `yaml:"mysql"`
}

type MySQLConfig struct {
	DSN       string `yaml:"dsn"`
	Host      string `yaml:"host"`
	Port      string `yaml:"port"`
	Username  string `yaml:"username"`
	Password  string `yaml:"password"`
	Database  string `yaml:"database"`
	Charset   string `yaml:"charset"`
	ParseTime bool   `yaml:"parse-time"`
	Loc       string `yaml:"loc"`
}

type QiniuKodoConfig struct {
	Path      string `yaml:"path"`
	Bucket    string `yaml:"bucket"`
	Domain    string `yaml:"domain"`
	AccessKey string `yaml:"access-key"`
	SecretKey string `yaml:"secret-key"`
}

func Load(path string) (Config, error) {
	var config Config
	content, err := os.ReadFile(path)
	if err != nil {
		return config, err
	}
	if err := yaml.Unmarshal(content, &config); err != nil {
		return config, err
	}
	return config, nil
}
