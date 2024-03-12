package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	Images  ImageConfig  `yaml:"images"`
	RootFs  RootFsConfig `yaml:"rootfs"`
	Version string       `yaml:"version"`
}
type ImageConfig struct {
	ImagePath string `yaml:"imagePath"`
}
type RootFsConfig struct {
	RootFsPath     string `yaml:"rootfsPath"`
	WorkLayerPath  string `yaml:"workLayerPath"`
	UpperLayerPath string `yaml:"upperLayerPath"`
	MntPath        string `yaml:"mntPath"`
	MntOldPath     string `yaml:"mntOldPath"`
}

var Cfg *Config

func Init() error {
	pwd, err := os.Getwd()
	if err != nil {
		fmt.Errorf("get pwd fail err=%s", err)
	}
	file, err := os.ReadFile(pwd + "/config/config.yaml")
	if err != nil {
		return err
	}
	err = yaml.Unmarshal(file, &Cfg)
	if err != nil {
		return err
	}
	return nil
}
