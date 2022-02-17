package config

type App struct {
	MaxItemsToSend   int `yaml:"maxItemsToSend"`
	SecondSleep      int `yaml:"secondSleep"`
	NumberOfAttempts int `yaml:"numberOfAttempts"`
	StartHour        int `yaml:"startHour"`
	StartMinute      int `yaml:"startMinute"`
	StartSecond      int `yaml:"startSecond"`
	SleepAfterSync   int `yaml:"sleepAfterSync"`
}
