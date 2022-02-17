package config

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
)

type Config struct {
	Store      Store      `yaml:"store"`
	Papi       Papi       `yaml:"papi"`
	App        App        `yaml:"app"`
	Prometheus Prometheus `yaml:"prometheus"`
}

func ReadConfig(filename string) (*Config, error) {
	config, err := readYaml(filename)
	if err != nil {
		config, err = readEnv()
		if err != nil {
			return nil, err
		}
	}

	return config, nil
}

func readEnv() (*Config, error) {
	var config Config

	config.Store.User = os.Getenv("DB_USER")
	config.Store.Pwd = os.Getenv("DB_PWD")
	config.Store.Dsn = os.Getenv("DB_DSN")
	config.Store.Port = os.Getenv("DB_PORT")
	config.Store.Database = os.Getenv("DB_DB")

	config.Papi.BaseUrl = os.Getenv("BASE_URL")
	config.Papi.Url = os.Getenv("URL")

	config.App.MaxItemsToSend, _ = strconv.Atoi(os.Getenv("MAX_ITEMS_TO_SEND"))

	config.App.SecondSleep, _ = strconv.Atoi(os.Getenv("SECOND_SLEEP"))
	config.App.NumberOfAttempts, _ = strconv.Atoi(os.Getenv("NUMBER_OF_ATTEMPTS"))
	config.App.StartHour, _ = strconv.Atoi(os.Getenv("START_HOUR"))
	config.App.StartMinute, _ = strconv.Atoi(os.Getenv("START_MINUTE"))
	config.App.StartSecond, _ = strconv.Atoi(os.Getenv("START_SECOND"))
	config.App.SleepAfterSync, _ = strconv.Atoi(os.Getenv("SLEEP_AFTER_SYNC"))

	config.Prometheus.Addr = os.Getenv("PROMETHEUS_ADDR")

	return &config, nil
}

func readYaml(filename string) (*Config, error) {
	var c *Config

	yamlFile, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to get config file: %v ", err)
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %v", err)
	}

	return c, nil
}
