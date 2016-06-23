package main

import (
	"flag"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os"
	fp "path/filepath"
	"strings"
)

type Config struct {
	DevKey  string `yaml:"devkey"`
	UserKey string `yaml:"userkey"`
}

type Metadata struct {
	Filetype string
	Filename string
}

func main() {
	// setting up flags
	config := flag.String("conf", "", "config file path")

	flag.Parse()

	if *config == "" {
		homeDir := os.Getenv("HOME")
		*config = strings.Join([]string{homeDir, ".pastebin.yaml"}, "/")
	}
}

// Read file to be pasted into memory
func LoadFile(filepath string) string {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Fatal("Filepath at '%s' does not exist", filepath)
	}

	fileContents, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal("Failed to read file at '%s'", filepath)
	}

	return string(fileContents)
}

func LoadFileMetadata(filepath string) Metadata {
	// need to remove the leading '.' from extension
	fileType := string(fp.Ext(filepath)[1:])

	// do we have an empty extension?
	if fileType == "" {
		fileType = "txt"
	}

	// getting filename
	fileName := strings.TrimRight(fp.Base(filepath), fileType)

	return Metadata{
		Filetype: fileType,
		Filename: fileName,
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
