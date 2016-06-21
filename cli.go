package main

import "fmt"
import "flag"
import yaml "gopkg.in/yaml.v2"

func main() {
	// setting up flags
	config := flag.String("conf", "", "config file path")

	flag.Parse()

	if *config == "" {
		homeDir := os.Getenv("HOME")
		*conf = strings.Join([]string{homedir, ".pastebin.yaml"})
	}
}
