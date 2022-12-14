# Configuration

Configuration module based on [configor](https://github.com/jinzhu/configor) tool.

## How to use

Initial config that will store configuration information should be placed in root folder named `app.default.json`. Write config structure of your app.

Example of config structure:

```go
type ConfigType struct {
 AppName   string
 Version   string
 Prefork   bool
}
```

Initialize and use config:

```go
config := config.Init[ConfigType]()
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

You can remove subscriber by given name on the fly as well:

```go
c.RemoveSubscriber("name_of_subscriber")
```

Library also support optional parameters with high order functions:

```go
config := config.Init[ConfigType](WithPath("./configuration_dir"), WithName("configuration_name"))
```
