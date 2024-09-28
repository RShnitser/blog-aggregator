package config

import(
	"os"
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	DbURL           string `json:"db_url"`
	CurrentUserName *string `json:"current_user_name"`
}

func(config *Config) SetUser(name *string)error{
	config.CurrentUserName = name

	home, err := os.UserHomeDir()
	if err != nil{
		return err
	}

	file, err := os.Open(home + "/.gatorconfig.json")
	if err != nil{
		return err
	}
	defer file.Close()

	json, err := json.Marshal(config)
	if err != nil{
		return err
	}

	file.Write(json)

	return nil

}

func Read()(Config, error){

	home, err := os.UserHomeDir()
	if err != nil{
		return Config{}, err
	}

	file, err := os.Open(home + "/.gatorconfig.json")
	if err != nil{
		return Config{}, err
	}

	defer file.Close()

	bytes, _ := ioutil.ReadAll(file)

	var cfg Config
	json.Unmarshal(bytes, &cfg)

	return cfg, nil
}