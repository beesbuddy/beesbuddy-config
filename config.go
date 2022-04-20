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
	File        string
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

	for i := 0; i < numberOfSubs; i++ {
		c.Subscribers = append(c.Subscribers, make(chan bool))
	}

	c.File = defaultConfig
	c.UpdateVersion(DONT_NOTIFY)
	if _, err := os.Stat(activeConfig); err == nil {
		c.File = activeConfig
	}

	err := configor.Load(&c.Cfg, c.File)
	if err != nil {
		log.Fatal("Configuration error: ", err)
	}

	c.File = activeConfig
	c.PersistToFile()

	return c
}

func (c *Data[T]) UpdateConfig(newConfig T) {
	c.Cfg = newConfig
	c.PersistToFile()
	c.UpdateVersion(SHOULD_NOTIFY)
}

func (c *Data[T]) UpdateVersion(n Notify) {
	c.Timestamp = strconv.FormatInt(time.Now().Unix(), 10)
	if n == SHOULD_NOTIFY {
		for i := 0; i < len(c.Subscribers); i++ {
			c.Subscribers[i] <- true
		}
	}
}

func (c *Data[T]) PersistToFile() {
	file, _ := json.MarshalIndent(c.Cfg, "", "	")
	_ = ioutil.WriteFile(c.File, file, 0644)
}
