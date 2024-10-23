package config

import (
	"encoding/json"
	"errors"
	"io"
	"os"
)

const ConfigFilename = ".gatorconfig.json"

type Config struct {
	dbURL	string
}

func Read(filePath string) (Config, error) {
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	dat, err := io.ReadAll(jsonFile)
	if err != nil {
		return Config{}, err
	}
	config := Config{}
	err = json.Unmarshal(dat, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func SetUser(config Config) error {
	return nil
}

func write(cfg Config) error {
	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := home + ConfigFilename
	return path, nil
}

