# Configuration

Configuration module based on [configor](https://github.com/jinzhu/configor) tool.

## How to use

Initial config that will store configuration information should be placed in root folder named `app_config.initial.json`. Write config structure of your app.

Example of config structure:

```go
type ConfigType struct {
 AppName   string `default:"worker"`
 Version   string `default:"1"`
 Prefork   bool   `default:"false"`
}
```

Initialize and use config:

```go
config := config.NewConfig[ConfigType]()
// access current configuration attributes
cfg := config.GetCfg()
cfg.AppName = "NewName"
// update current configuration
config.UpdateConfig(cfg)
```

If you have modules which needs to be notified on config change, add a listener/subscriber:

```go
c.AddSubscriber("name_of_subscriber")
```

Implement waiting goroutine for config change on the fly in your modules:

```go
_ = <-config.GetSubscriber("name_of_subscriber")
```
