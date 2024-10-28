package config

import (
	"encoding/json"
	"io"
	"os"
)

const configFilename = ".gatorconfig.json"

type State struct {
	Config 		*Config 	
}

type Config struct {
	DBURL				string
	CurrentUsername		string
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	jsonFile, err := os.Open(filePath)
	if err != nil {
		return Config{}, err
	}
	defer jsonFile.Close()

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

func (cfg *Config)SetUser(username string) error {
	cfg.CurrentUsername = username
	err := cfg.write()
	if err != nil {
		return err
	}
	return nil
}

func (cfg *Config)write() error {
	dat, err := json.Marshal(&cfg)
	if err != nil {
		return err
	}

	filePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	os.WriteFile(filePath, dat, 0644)

	return nil
}

func getConfigFilePath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	path := home + "/" + configFilename
	return path, nil
}

