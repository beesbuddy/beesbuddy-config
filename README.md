# config
Configuration module for `tableassistant` based on [configor](github.com/jinzhu/configor) tool.

### How to use
Write config structure of your app. Actual config should be placed in root folder named `appConfig.json` 
Example:
```
type ConfigType struct {
	AppName   string `default:"worker"`
	Version   string `default:"1"`
	Prefork   bool   `default:"false"`
}
```

If you have modules which needs to be notified on config change, create similar enum:
```
type ConfigSubscriber int

const (
	ONE_SUB ConfigSubscriber = iota
	SECOND_SUB
	NUMBER_OF_SUBS
)
```

Initialize and use config:
```
config := config.NewConfig[ConfigType](int(NUMBER_OF_SUBS))

config.Cfg // access to current config
config.UpdateConfig({newConfig}) // update current config on the fly
```

Implement waiting goroutine for config change in your modules:
```
_ = <-config.Subscribers[SECOND_SUB]
```



