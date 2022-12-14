# Configuration

Configuration module based on [configor](https://github.com/jinzhu/configor) tool.

## How to use

Write config structure of your app. Actual config should be placed in root folder named `app_config.default.json`
Example:

```go
type ConfigType struct {
 AppName   string `default:"worker"`
 Version   string `default:"1"`
 Prefork   bool   `default:"false"`
}
```

If you have modules which needs to be notified on config change, create similar enum:

```go
type ConfigSubscriber int
const (
 FIRST_SUB ConfigSubscriber = iota
 SECOND_SUB
 NUMBER_OF_SUBS
)
```

Initialize and use config:

```go
config := config.NewConfig[ConfigType](int(NUMBER_OF_SUBS))

cfg := config.Cfg // access current config attributes
cfg.AppName = "NewName"
config.UpdateConfig(cfg) // update current config on the fly
```

Implement waiting goroutine for config change on the fly in your modules:

```go
_ = <-config.Subscribers[SECOND_SUB]

```
