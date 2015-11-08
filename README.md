# Go Config
The `go-config` package is used to simplify configuration for go applications. 

# Installation
To install this package, run:
```
go get github.com/zpatrick/go-config
```

# The Basics
The `go-config` package has three main components: **providers**, **settings**, and the **config** object. 
* **Providers** load settings for your application. This could be a file, environment variables, or some other means of configuration.
* **Settings** represent the configuration options for your application. Settings are represented as key/value pairs. 
* **Config** holds all of the providers in your application. This object holds convenience functions for loading, retrieving, and converting your settings.

# Built In Providers 
* `INIFile` - Loads settings from a `.ini` file
* `JSONFile`  - Loads settings from a `.json` file
* `YAMLFile` - Loads settings from a `.yaml` file
* `Environment` - Loads settings from environment variables 
* `Static` - Loads settings from an in-memory map

## Single Provider Example
Most application use a single file for configuration. 
In this example, the application uses the file `config.ini` with the following contents: 
```
[global]
timeout=30
frequency=0.5

[local]
time_zone=PST
enabled=true
```
First, create a **Provider** that will be responsible for loading configuration values from the `config.ini` file.
Since `config.ini` is in  [INI](https://code.google.com/p/minini/wiki/INI_File_Syntax) format, use the `INIFile` provider:
```
    iniFile := config.NewINIFile("config.ini")
```

Next, create a **Config** object to manage the application's providers. Since this application only has one provider, `iniFile`, we pass a list with only one object into the constructor. 
```
    c := config.NewConfig([]config.Provider{iniFile})
```

Althought it isn't required, it is a good idea to call the `Load()` function on your config object after creating it. This will notify you if there was an error loading the providers.
```
    if err := c.Load(); err != nil{
        log.Fatal(err)
    }
```
Calling the `Load()` function is not required. Each time a setting is requested, the `Load()` function is called first. If you are concerned about performance of calling `Load()` so frequently, see the [Advanced](#Advanced) section.  

The `Config` object can be used to lookup settings from your providers. When performing a lookup on a setting, the key is a period-delimited string. For `ini` files, this means lookups are performed with `<section>.<item>` keys. For example, the settings in `config.json` can be looked up with the following keys:
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
All of the settings in your application are stored as strings. As shown above, the `Config` object has convenience functions for type conversions: `Int(), Float(), Bool()`. 

## Validation
The `Config` object has an optional `Validate` function that is called after the settings are loaded from the providers. Taking the `config.ini` example from above, we can create a `Validate` function that makes sure the `global.timeout` setting exists: 
```
   iniFile := config.NewINIFile("config.ini")
   c := config.NewConfig([]config.Provider{iniFile})
   
    c.Validate = func(settings map[string]string) error {
        if val, ok := settings["global.timeout"]; !ok {
            return fmt.Errorf("Required setting 'global.timeout' not set!")
        }
    }
```

## Optional Settings 
Looking up settings using the `String()`, `Int()`, `Float()`, or `Bool()` functions will error if the setting does not exist. 
The functions `StringOr()`, `IntOr()`, `FloatOr()`, and `BoolOr()` can be used to provide a default if the setting does not exist:
```
    timeout, err := c.IntOr("global.timeout", 30)
    ...
    frequency, err := c.FloatOr("global.frequency", 0.5)
    ...
    tz, err := c.StringOr("local.time_zone", "PST")
    ...
    enabled, err := c.BoolOr("local.enabled", true)
    ...
```

## Get All Settings 
All of the setting can be retrieved using the `Settings()` function:
```
    settings, err := c.Settings()
    if err != nil{
        return err
    }
    
    for key, val := range settings {
        fmt.Printf("%s = %s \n", key, val)
    }
```

## Multiple Providers Example
Many applications have more than one means of configuration. 
Building off of the [Single Provider Example](#Single_Provider_Example) above, we will add environment variables as a means of configuration in addition to the `config.ini` file. We will let environment variables override settings in `config.go`. 

As before, create a Provider that will be responsible for loading configuration values from the `config.ini` file:
```
    iniFile := config.NewINIFile("config.ini")
```

Next, create an `Environment` provider that will be responsible for loading configuration values fron environment variables. The `Environment` provider takes a map that associates setting keys with environment variables.
Since we want the environment variables to override the same settings in `config.ini`, construct the map like so:
```
    mappings := map[string]string{
        "global.timeout": "APP_TIMEOUT",
        "globabl.frequency": "APP_FREQUENCY",
        "local.time_zone": "APP_TIMEZONE",
        "local.enabled": "APP_ENABLED",
    }
    
    env := config.NewEnvironment(mappings)
```
Since we have two providers, we add both of them to the `Config` object. The position of the provider in the list determines the ordering of settings lookups. Since we want environment variables to override the values in `config.ini`, we put the `Environment` provider later in the list:
```
    providers := []config.Providers{iniFile, env}
    c := config.NewConfig(providers)
```

# Advanced

## Defaults
Managing default values for settings can be accomplished multiple ways:
* A configuration file that contains all of the defaults
* Using `Or` functions with defaults (e.g. `StringOr(...)`, `IntOr(...)`, etc.)
* Using the `Static` provider

The `Static` provider takes key value mappings for settings and simply returns those values when `Load()` is called. This is a nice pattern as it doesn't require additional configuration files and it places all defaults into a single place in your code. Make sure to set your defaults as the first provider in your application so they can be overridden by other providers:
```
    mappings := map[string]string{
        "global.timeout": "30",
        "global.enabled": "true",
    }
    
    defaults := config.NewStatic(mappings)
    iniFile := config.NewINIFile("config.ini")
    
    providers := []config.Provider{defaults, iniFile}
    c := config.NewConfig(providers)
```

## Loading Patterns
Each time a lookup is performed, the `Load()` function is called on each provider. This can lead to poor performance and be unecessary for certain providers. There are two built in objects which change how frequently loads are performed:
* `OnceLoader` - Loads the provider's settings one time
* `CachedLoader` - Loads the provider's settings at least one time and caches the results. The `Invalidate()` function can be called to force a new load next time a lookup is performed.

Using the [Multiple Providers Example](#Multiple_Providers_Example) above, we will use the `OnceLoader` for the `iniFile` provider and keep the default behavior for the `Environment` provider (perform a load each time a lookup is requested):
```
    env := config.NewEnvironment(...)
    iniFile := config.NewINIFile("config.ini")
    iniFileOnce := config.NewOnceLoader(iniFile)
    
    providers := []config.Provider{iniFileOnce, env}
    c := config.NewConfig(providers)
```

An example using the `CachedLoader` instead of `OnceLoader`:
```
    env := config.NewEnvironment(...)
    iniFile := config.NewINIFile("config.ini")
    iniFileCached := config.NewCachedLoader(iniFile)
    
    providers := []config.Provider{iniFileCached, env}
    c := config.NewConfig(providers)
    ...
    // will force iniFileCached to load next time a lookup is performed
    c.Invalidate()
    v, err := c.String("global.timeout")
    ...
```


## Custom Providers
Custom providers must fulfill the `Provider` interface:
```
    type Provider interface {
        Load() (map[string]string, error)
    }
```

The `Load()` function returns settings as key/value pairs. Providers flatten namespaces using period-delimited strings. 
For example, the following providers and content:

INI File:
```
[global]
timeout = 30
```
YAML File:
```
global:
    timeout: 30
```

JSON File:
```
{
    "global": { 
        "timeout": 30
    }
}
```

All resolve the `global.timeout` setting to `30`. 
This is allows providers to override and lookup settings using a command key format. 



