# Go Config
The `go-config` package is used to simplify configuration for go applications. 

# Installation
To install this package, run:
```
go get github.com/zpatrick/go-config
```

# Getting Started
You start by creating a **Provider**. Providers load configuration settings for your application. 
For this example, assume the file `config.ini` holds the configuration for the application with the following format:
```
[global]
timeout=30
frequency=0.5

[local]
time_zone=pst
enabled=true
```

First, create a provider for the file:
```
    ini := config.NewINIFile("config.ini")
```

Now, create a new `Config` object with the `ini` provider. The `Config` object will use all of the `providers` to lookup configuration settings.
```
c := config.NewConfig([]config.Provider{ini})
```

Now you can lookup configuration settings and their types as required:
```
    timeout, err := c.Int("global.timeout")
    ...
    frequency, err := c.Float("global.frequency")
    ...
    tz, err := c.String("local.time_zone")
    ...
    enabled, err := c.Bool("local.enabled")
    ...
```

### Defaults
Check if a setting exists 

# Multiple Providers
TODO: Config file + environment variables

In the following example, the `timeout` variable will resolve to `60`. 
```
    ini :=  config.NewINIFile("global_config.ini")
    environment := config.NewEnvironment(...)

    c := config.NewConfig([]config.Provider{ini, environment})
    
    timeout, err := c.String("server.timeout")
    ...
```




# Defaults
Create a static provider or use StringOr, IntOr, BoolOr
# Custom Providers
Fulfill the `Provider` interface. 