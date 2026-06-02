package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	QiniuKodo QiniuKodoConfig `yaml:"qiniu-kodo"`
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
