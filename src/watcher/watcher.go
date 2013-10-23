package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Watch string              `json:"watch"`
	Files map[string][]string `json:"files"`
}

var configPath string

func init() {
	flag.StringVar(&configPath, "config", "config.json", "path to configuration json file")
	flag.Parse()
}

func loadConfig() *Config {
	var watcherConfig Config

	configFile, err := os.Open(configPath)
	defer configFile.Close()

	if err != nil {
		log.Fatal("could not open config file", configPath, err)
	}
	content, err := ioutil.ReadAll(configFile)
	if err != nil {
		log.Fatal("could not read config file", configPath, err)
	}
	err = json.Unmarshal(content, &watcherConfig)
	if err != nil {
		log.Fatal("could not parse the contents of the config file", err)
	}

	return &watcherConfig
}

func watchFile(path string) chan bool {
	var ch chan bool = make(chan bool)
	var lastTime int64
	lastTime = 0

	go func() {
		for {
			if fp, err := os.Open(path); err == nil {
				fileInfo, err := fp.Stat()
				if err == nil {
					if fileInfo.ModTime().Unix() > lastTime {
						lastTime = fileInfo.ModTime().Unix()
						ch <- true
					}
				} else {
					log.Println("Could not get file info for", path, err)
				}
				fp.Close()
			} else {
				log.Println("Could not open file ", path, err)
			}
			time.Sleep(time.Second * 10)
		}
	}()

	return ch
}

func (config *Config) doUploads() {
	for file, urls := range config.Files {
		for _, url := range urls {
			if fp, err := os.Open(file); err != nil {
				log.Println("could not open upload file", file, err)
			} else {
				_, err := http.Post(url, "text/plain", fp)
				fp.Close()
				if err != nil {
					log.Println("error posting file", file, err)
				}
			}
		}
	}
}

func (config *Config) watch() {
	ch := watchFile(config.Watch)
	for {
		<-ch
		config.doUploads()
	}
}

func main() {
	config := loadConfig()
	config.watch()

	for {
		time.Sleep(time.Second)
	}
}
