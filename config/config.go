package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
	"os/user"
	"path/filepath"
)

type Configuration struct {
	BindPort string `json:BindPort`
	Debug    bool   `json:Debug`
}

// Initialise default configuration.
var Config *Configuration = &Configuration{BindPort: "6483", Debug: true}

func GetDataDir() string {

	user, err := user.Current()
	if err != nil {
		log.Fatalf(err.Error())
	}

	configdir := filepath.Join(user.HomeDir, ".config", "odbauth")

	if _, err := os.Stat(configdir); os.IsNotExist(err) {
		// config dir not exist
		err = os.Mkdir(filepath.Join(user.HomeDir, ".config", "odbauth"), 0755)

		if err != nil {
			panic(err)
		}
	}

	return configdir
}

func InitialiseConfig() *Configuration {
	defer func() {
		if err := recover(); err != nil {
			log.Println("Fatal exception when parsing the configuration data. CHECK YOUR CONFIG!", err)
		}
	}()

	config_data_location := filepath.Join(GetDataDir(), "config.json")

	if _, err := os.Stat(config_data_location); os.IsNotExist(err) {
		// config dir not exist

		log.Fatalln("Configuration doesn't exist! Please place your config in ~/.config/odbauth/config.json in order to start the site.")
	}

	if data, err := os.ReadFile(config_data_location); err == nil {
		// got file.
		var parsed Configuration

		// parse. panic on error.
		if err := json.Unmarshal(data, &parsed); err != nil {
			panic(err)
		}

		// return the parsed json.
		return &parsed
	} else if errors.Is(err, os.ErrNotExist) {
		// does not exist.
		log.Fatalln("Configuration doesn't exist! Please place your config in ~/.config/odbauth/config.json in order to start the site.")
	} else {
		// schrodingers file. presume non-existant.
		log.Fatalln("Configuration doesn't exist! Please place your config in ~/.config/odbauth/config.json in order to start the site.")
	}

	return Config
}
