package config

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/jinzhu/configor"
)

var (
	defaultConfig = "appConfig.json"
	activeConfig  = "appConfig_active.json"
)

type Data[T any] struct {
	file        string
	Timestamp   string
	Cfg         T
	Subscribers []chan bool
}

type Notify int

const (
	SHOULD_NOTIFY Notify = iota
	DONT_NOTIFY
)

func NewConfig[T any](numberOfSubs int) *Data[T] {
	c := new(Data[T])

	c.file = defaultConfig
	if _, err := os.Stat(activeConfig); err == nil {
		c.file = activeConfig
	}
	err := configor.Load(&c.Cfg, c.file)
	if err != nil {
		log.Fatal("Configuration error: ", err)
	}
	c.file = activeConfig

	for i := 0; i < numberOfSubs; i++ {
		c.Subscribers = append(c.Subscribers, make(chan bool))
	}
	c.updateVersion(DONT_NOTIFY)
	c.persistToFile()

	return c
}

func (c *Data[T]) UpdateConfig(newConfig T) {
	c.Cfg = newConfig
	c.persistToFile()
	c.updateVersion(SHOULD_NOTIFY)
}

func (c *Data[T]) updateVersion(n Notify) {
	c.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	if n == SHOULD_NOTIFY {
		for i := 0; i < len(c.Subscribers); i++ {
			c.Subscribers[i] <- true
		}
	}
}

func (c *Data[T]) persistToFile() {
	file, _ := json.MarshalIndent(c.Cfg, "", "	")
	_ = ioutil.WriteFile(c.file, file, 0644)
}
