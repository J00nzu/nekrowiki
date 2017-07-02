package main

import (
	"time"
	"fmt"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"os"
	"log"
	"bufio"
)

type configValues struct {
	HomePage                  string
	SessionExpirationTime     time.Duration
	SessionStoreCleanInterval time.Duration
	AllowedQuestAccess        bool
}

var config = DefaultConfig()

func DefaultConfig() configValues {
	fmt.Print("Initializing configs\n")

	values := configValues{}
	values.HomePage = "/home"
	values.SessionExpirationTime = time.Hour
	values.SessionStoreCleanInterval = time.Hour
	values.AllowedQuestAccess = false;

	return values
}

func LoadConfiguration(cfgFile string) {
	var conf configValues

	tomlData, err := ioutil.ReadFile(cfgFile)

	log.Println("Loading Config")

	if string(tomlData) == "" {
		SaveConfiguration(cfgFile)
		return
	}

	if err != nil {
		log.Println(err)

		SaveConfiguration(cfgFile)
		return
	}

	if _, err := toml.Decode(string(tomlData), &conf); err != nil {
		log.Println(err)
		SaveConfiguration(cfgFile);
	} else {
		config = conf
	}

}

func SaveConfiguration(cfgFile string) {

	log.Println("Saving Config")

	file, err := os.Create(cfgFile)
	if (err != nil) {
		log.Println(err)
		panic(err)
	}

	defer file.Close()

	w := bufio.NewWriter(file)
	enc := toml.NewEncoder(w)
	err = enc.Encode(config)

	if (err != nil) {
		log.Println(err)
		panic(err)
	}
}


