package config

import "github.com/tkanos/gonfig"

//Struct
type Configuration struct {
	DB_Username string
	DB_PASSWORD string
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
}

func GetConfig() Configuration {
	conf := Configuration{}
	err := gonfig.GetConf("config/config.json", &conf)

	if err != nil {
		return conf
	}

	return conf
}
