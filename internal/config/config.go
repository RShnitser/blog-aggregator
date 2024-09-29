package config

import(
	"os"
	"encoding/json"
	"path/filepath"
)


type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func getConfigFilePath() (string, error) {
	const configFileName = "/.gatorconfig.json"
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	fullPath := filepath.Join(home, configFileName)
	return fullPath, nil
}


func(config *Config) SetUser(name string)error{
	config.CurrentUserName = name

	path, err := getConfigFilePath()
	if err != nil{
		return err
	}

	file, err := os.Create(path)
	if err != nil{
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(config)
	if err != nil{
		return err
	}

	return nil
}

func Read()(Config, error){

	path, err := getConfigFilePath()
	if err != nil{
		return Config{}, err
	}

	file, err := os.Open(path)
	if err != nil{
		return Config{}, err
	}

	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := Config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		return Config{}, err
	}

	return cfg, nil
}