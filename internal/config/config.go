package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

var BotToken = os.Getenv("TELEGRAM_BOT_TOKEN")

const GroupChatID int64 = -865707097

type ProviderData struct {
	Name         string `yaml:"name"`
	LocalAddress string `yaml:"local_address"`
	WAN          string `yaml:"wan"`
	Saturation   int64  `yaml:"saturation"`
}

type ProviderConfig struct {
	Providers []ProviderData `yaml:"providers"`
}

type BtsData struct {
	Name         string `yaml:"name"`
	LocalAddress string `yaml:"local_address"`
	WAN          string `yaml:"wan"`
}

type BtsConfig struct {
	Bts []BtsData `yaml:"bts"`
}

func LoadProviderConfig() (ProviderConfig, error) {
	log.Println("Reading provider configuration file")
	filename, _ := filepath.Abs("../provider.config.yml")
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return ProviderConfig{}, err
	}

	var config ProviderConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return ProviderConfig{}, err
	}

	return config, nil
}

func LoadBtsConfig() (BtsConfig, error) {
	log.Println("Reading bts configuration file")
	filename, _ := filepath.Abs("../bts.config.yml")
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		return BtsConfig{}, err
	}

	var config BtsConfig
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		return BtsConfig{}, err
	}

	return config, nil
}
