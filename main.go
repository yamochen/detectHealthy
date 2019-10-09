package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"gopkg.in/ini.v1"
)

func log(msg string) {
	logFile, _ := os.OpenFile("./home/logs/detechHealthy.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	info := fmt.Sprintf("[%s] %s", time.Now().Format("2006-01-02 15:04:05.99999"), msg)
	logFile.WriteString(info)
}

func main() {
	config, err := loadConfig()
	if err != nil {
		panic(err)
	}

	resp, err := http.Get(config.URL)
	if err != nil {
		msg := fmt.Sprintf("ERROR:PID:%d , ERROR-Can't connect to %s\n", os.Getpid(), config.URL)
		log(msg)
		return
	}
	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	reader := strings.NewReader(string(data))
	j := json.NewDecoder(reader)
	j.UseNumber()

	status := make(map[string]string)

	if err := j.Decode(&status); err != nil {
		panic(err)
	}

	msg := fmt.Sprintf("%s:PID:%d ,response: %s\n", status["mysql_status"], os.Getpid(), string(data))
	log(msg)
}

type Config struct {
	Logger int    `ini:"logger"`
	URL    string `ini:"URL"`
}

func loadConfig() (config Config, err error) {
	cfg, err := ini.Load("./config.ini") //load config
	if err != nil {
		return config, err
	}

	err = cfg.MapTo(&config) //Parser To Struct
	if err != nil {
		return config, err
	}

	return config, nil
}
