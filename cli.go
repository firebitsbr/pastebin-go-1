package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	yaml "gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	fp "path/filepath"
	"strings"
)

/************************************
 ********** STRUCT DEFS *************
 ************************************/
type Config struct {
	DevKey  string `yaml:"devkey"`
	UserKey string `yaml:"userkey"`
}

type Metadata struct {
	Filetype     string
	Filename     string
	Filecontents string
}

type PostData struct {
	DevKey     string `json:"api_dev_key"`
	Option     string `json:"api_option"`
	Code       string `json:"api_paste_code"`
	UserKey    string `json:"api_user_key"`
	Name       string `json:"api_paste_name"`
	Filetype   string `json:"api_paste_format"`
	Privacy    int    `json:"api_paste_private"`
	Expiration string `json:"api_paste_expire_date"`
}

/************************************
 *********** CONSTANTS **************
 ************************************/
const baseUrl string = "http://pastebin.com/api/api_post.php"

/************************************
 *************** MAIN ***************
 ************************************/
func main() {
	// setting up flags
	config := flag.String("conf", "", "config file path")
	expiration := flag.String("exp", "", "paste expiration date")
	privacy := flag.Int("priv", -1, "post privacy settings")

	flag.Parse()

	if *config == "" {
		homeDir := os.Getenv("HOME")
		*config = strings.Join([]string{homeDir, ".pastebin.yaml"}, "/")
	}

	if *privacy != 0 && *privacy != 1 && *privacy != 2 && *privacy != -1 {
		log.Fatal("Privacy should be 0, 1, or 2 (defaults to 0)")
	}

	// currently only supporting PB of a single file
	tail := flag.Args()

	if len(tail) != 1 {
		log.Fatal("Expecting one file as input")
	}

	fileMeta := LoadFile(tail[0])
	pbConf := LoadConfig(*config)

	pbUrl := GeneratePaste(fileMeta, pbConf, *expiration, *privacy)

	fmt.Printf("Got back url of '%s'\n", pbUrl)
}

// Load file to be read into memory.
// @Params: filepath - path where file is is located
// @Return: Metadata struct
func LoadFile(filepath string) Metadata {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		log.Fatal("Filepath at '%s' does not exist", filepath)
	}

	fileContents, err := ioutil.ReadFile(filepath)
	if err != nil {
		log.Fatal("Failed to read file at '%s'", filepath)
	}

	// need to remove the leading '.' from extension
	fileType := string(fp.Ext(filepath)[1:])

	// do we have an empty extension?
	if fileType == "" {
		fileType = "txt"
	}

	// getting filename
	fileName := strings.TrimRight(fp.Base(filepath), fileType)

	return Metadata{
		Filetype:     fileType,
		Filename:     fileName,
		Filecontents: string(fileContents),
	}
}

// Load configuration into memory
// @Params: confpath - filepath of configuration file
// @Return: Config struct
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

// Generate POST to pasteboard API
// @Params: meta - file Metadata struct
//				 - conf Config struct
//               - expiration date string
//               - privacy denoting visibility of paste
// @Return: url provided by pastebin
func GeneratePaste(meta Metadata, conf Config, expiration string, privacy int) string {
	// privacy defaults to public
	if privacy == -1 {
		privacy = 0
	}

	if expiration == "" {
		expiration = "N"
	}

	// load this into struct
	data := PostData{
		DevKey:     conf.DevKey,
		Code:       meta.Filecontents,
		Privacy:    privacy,
		Name:       meta.Filename,
		Expiration: expiration,
		Filetype:   meta.Filetype,
		UserKey:    conf.UserKey,
		Option:     "paste",
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Fatal("Could not marshal request into json")
	}

	req, err := http.NewRequest("POST", baseUrl, bytes.NewBuffer(jsonData))
	if err != nil {
		log.Fatal("Error creating HTTP request")
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Error sending HTTP request")
	}

	defer resp.Body.Close()

	log.Printf("Got response with %s code\n", resp.Status)
	body, _ := ioutil.ReadAll(resp.Body)

	return string(body)
}
