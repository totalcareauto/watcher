package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"time"

	"watcher/Godeps/_workspace/src/github.com/stvp/rollbar"
)

type Config struct {
	Production   bool                `json:"production"`
	LogFile      string              `json:"log_file""`
	RollbarToken string              `json:"rollbar_token"`
	Watch        string              `json:"watch"`
	Files        map[string][]string `json:"files"`
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

func checkForStale(path string) {
	go func() {
		for {
			if time.Now().Hour() >= 7 {
				if fp, err := os.Open(path); err == nil {
					if fileInfo, err := fp.Stat(); err == nil {
						if math.Abs(time.Now().Sub(fileInfo.ModTime()).Hours()) > 1 {
							message := "Watch file " + fileInfo.Name() + " is over one hour old."
							logger.Println(message)
							err = errors.New(message)
							rollbar.Error("error", err)
						}
					}
					fp.Close()
				}
			}

			time.Sleep(time.Minute * 10)
		}
	}()
}

func watchFile(path string) chan bool {
	var ch chan bool = make(chan bool)
	var lastTime int64
	lastTime = 0

	checkForStale(path)

	go func() {
		for {
			if fp, err := os.Open(path); err == nil {
				if fileInfo, err := fp.Stat(); err == nil {
					if fileInfo.ModTime().Unix() > lastTime {
						lastTime = fileInfo.ModTime().Unix()
						ch <- true
					}
				} else {
					logger.Println("Could not get file info for", path, err)
					rollbar.Error("error", err)
				}
				fp.Close()
			} else {
				logger.Println("Could not open file ", path, err)
				rollbar.Error("error", err)
			}
			time.Sleep(time.Second * 10)
		}
	}()

	return ch
}

func doUpload(file, url string) {
	defer func() {
		if errStr := recover(); errStr != nil {
			fullErrStr := fmt.Sprintln("Error (", errStr, ") uploading file ", file, " to URL ", url)
			logger.Println(fullErrStr)
			rollbar.Error("error", errors.New(fullErrStr))
		}
	}()
	if fp, err := os.Open(file); err != nil {
		log.Println("Could not open upload file", file, err)
		rollbar.Error("error", err)
	} else {
		defer fp.Close()
		client := &http.Client{
			Timeout: time.Minute * 3,
		}
		resp, err := client.Post(url, "text/plain", fp)
		if err != nil {
			logger.Println("Error uploading file", file, err)
			rollbar.Error("error", err)
		} else {
			logger.Println("Uploaded file ", file)
			resp.Body.Close()
		}
	}
}

func (config *Config) doUploads() {
	defer func() {
		if errStr := recover(); errStr != nil {
			log.Println("error in doUploads(): ", errStr)
		}
	}()
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
		logger.Println("Watched file changed.  Attempting uploads")
		config.doUploads()
	}
}

func (config *Config) initialize() {
	// setup rollbar
	rollbar.Token = config.RollbarToken
	if config.Production {
		rollbar.Environment = "production"
	} else {
		rollbar.Environment = "development"
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
