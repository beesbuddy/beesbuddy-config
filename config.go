package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"time"
)

var (
	defaultConfig = "%s.defualt.json"
	activeConfig  = "%s.json"
)

type config[T any] struct {
	activeFile  string
	cfg         T
	subscribers map[string](chan bool)
	timestamp   string
}

const (
	MARSHAL_INDENT     = "	"
	EMPTY_SPACE        = ""
	RW_RW_R_PERMISSION = 0664
)

func Init[T any](name string) (*config[T], error) {
	c := &config[T]{}
	c.subscribers = make(map[string]chan bool)
	activeConfigFilename := fmt.Sprintf(activeConfig, name)
	defaultConfigFilename := fmt.Sprintf(defaultConfig, name)
	activeFileExists := fileExists(activeConfigFilename)
	defaultFileExists := fileExists(defaultConfigFilename)

	if activeFileExists {
		c.activeFile = activeConfigFilename
	} else if defaultFileExists {
		c.activeFile = defaultConfigFilename
	} else {
		return nil, fmt.Errorf("no configuration files found")
	}

	err := c.load()
	if err != nil {
		return nil, fmt.Errorf("failed at load from file: %v", err)
	}

	c.updateTimestamp()

	if !activeFileExists {
		c.activeFile = activeConfigFilename
		err = c.persist()
		if err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *config[T]) updateTimestamp() {
	c.timestamp = strconv.FormatInt(time.Now().Unix(), 10)
}

func (c *config[T]) Update(newConfig T) error {
	c.cfg = newConfig

	err := c.persist()

	if err != nil {
		return err
	}

	c.updateTimestamp()

	for _, channel := range c.subscribers {
		// Do not notify subscriber through channel if it was aleady notified
		if len(channel) != 0 {
			continue
		}

		channel <- true
	}

	return nil
}

func (c *config[T]) persist() error {
	file, err := json.MarshalIndent(c.cfg, EMPTY_SPACE, MARSHAL_INDENT)

	if err != nil {
		return fmt.Errorf("failed at marshal json: %v", err)
	}

	err = os.WriteFile(c.activeFile, file, RW_RW_R_PERMISSION)

	if err != nil {
		return fmt.Errorf("failed at write to file: %v", err)
	}

	return nil
}

func (c *config[T]) load() error {
	configFile, err := os.Open(c.activeFile)
	if err != nil {
		return err
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&c.cfg); err != nil {
		return err
	}

	return nil
}

func fileExists(filename string) bool {
	if _, err := os.Stat(filename); err == nil {
		return true
	}

	return false
}

func (c *config[T]) GetSubscriber(key string) chan bool {
	return c.subscribers[key]
}

func (c *config[T]) AddSubscriber(key string) {
	c.subscribers[key] = make(chan bool, 1)
}

func (c *config[T]) GetTimestamp() string {
	return c.timestamp
}

func (c *config[T]) GetSubscribers() map[string](chan bool) {
	return c.subscribers
}

func (c *config[T]) GetCfg() *T {
	return &c.cfg
}
