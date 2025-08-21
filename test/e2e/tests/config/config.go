package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type EnvConfig struct {
	PacGoServer          string `yaml:"PAC_GO_SERVER"`
	MongoDBURI           string `yaml:"MONGODB_URI"`
	KeycloakHost         string `yaml:"KEYCLOAK_HOSTNAME"`
	KeycloakRealm        string `yaml:"KEYCLOAK_REALM"`
	KeycloakClientID     string `yaml:"KEYCLOAK_CLIENT_ID"`
	KeycloakClientSecret string `yaml:"KEYCLOAK_CLIENT_SECRET"`
}

type ConfigFile struct {
	Test EnvConfig `yaml:"test"`
}

var Current EnvConfig

func LoadConfig(configPath, env string) error {
	if path := os.Getenv("CONFIG_PATH"); path != "" {
		configPath = path
	}

	str, _ := os.Getwd()
	fmt.Println("STR!!!!: ", str)

	f, err := os.Open(configPath)
	if err != nil {
		return fmt.Errorf("failed to open %s: %w", configPath, err)
	}
	defer f.Close()

	var cfg ConfigFile
	decoder := yaml.NewDecoder(f)
	if err := decoder.Decode(&cfg); err != nil {
		return fmt.Errorf("failed to decode config.yaml: %w", err)
	}

	switch env {
	case "test":
		Current = cfg.Test
	default:
		return fmt.Errorf("unknown environment: %s", env)
	}

	return nil
}
