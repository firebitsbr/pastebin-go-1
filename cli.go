package main

import "fmt"
import "flag"
import yaml "gopkg.in/yaml.v2"
import "os"
import "log"

type Config struct {
	DevKey  string `yaml:"devkey"`
	UserKey string `yaml:"userkey"`
}

func main() {
	// setting up flags
	config := flag.String("conf", "", "config file path")

	flag.Parse()

	if *config == "" {
		homeDir := os.Getenv("HOME")
		*conf = strings.Join([]string{homedir, ".pastebin.yaml"})
	}
}

func LoadConfig(confpath string) Config {
	if _, err := os.Stat(confpath); os.IsNotExist(err) {
		log.Fatal("Config at '%s' does not exist", confpath)
	}

	configContents, err := ioutil.ReadFile(confpath)
	if err != nil {
		log.Fatal("Could not read config at '%s'", confpath)
	}

	var config Config
	err = yaml.Unmarshal(configContents, &config)
	if err != nil {
		log.Fatal("Could not unmarshal configuration file")
	}
	return config
}
