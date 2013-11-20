package main

import (
	"encoding/json"
	"flag"
	"github.com/tobi/airbrake-go"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Config struct {
	Production  bool                `json:"production"`
	LogFile     string              `json:"log_file""`
	AirBrakeKey string              `json:"air_brake_key"`
	Watch       string              `json:"watch"`
	Files       map[string][]string `json:"files"`
}

var configPath string
var logger *log.Logger

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
					logger.Println("Could not get file info for", path, err)
					airbrake.Notify(err)
				}
				fp.Close()
			} else {
				logger.Println("Could not open file ", path, err)
				airbrake.Notify(err)
			}
			time.Sleep(time.Second * 10)
		}
	}()

	return ch
}

func doUpload(file, url string) {
	if fp, err := os.Open(file); err != nil {
		log.Println("Could not open upload file", file, err)
	} else {
		defer fp.Close()
		resp, err := http.Post(url, "text/plain", fp)
		resp.Body.Close()
		if err != nil {
			logger.Println("Error uploading file", file, err)
			airbrake.Notify(err)
		} else {
			logger.Println("Uploaded file ", file)
		}
	}
}

func (config *Config) doUploads() {
	for file, urls := range config.Files {
		for _, url := range urls {
			logger.Println("uploading ", file, " to ", url)
			doUpload(file, url)
		}
	}
}

func (config *Config) watch() {
	ch := watchFile(config.Watch)
	for {
		<-ch
		logger.Println("Watched file changed attempting uploads")
		config.doUploads()
	}
}

func (config *Config) initialize() {
	// setup airbrake
	airbrake.ApiKey = config.AirBrakeKey
	if config.Production {
		airbrake.Environment = "production"
	}

	// open the log output file
	if logFp, err := os.OpenFile(config.LogFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err == nil {
		logger = log.New(logFp, "", (log.Ldate | log.Ltime | log.Lmicroseconds | log.Lshortfile))
	} else {
		log.Fatal("Could not open log file ", config.LogFile)
	}
}

func main() {

	config := loadConfig()
	config.initialize()
	config.watch()

	for {
		time.Sleep(time.Second)
	}
}
