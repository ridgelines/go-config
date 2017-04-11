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
* `TOMLFile` - Loads settings from a `.toml` file
* `CLI` - Loads settings from a [codegansta/cli](https://github.com/codegangsta/cli) context
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
Since `config.ini` is in `ini` format, use the `INIFile` provider:
```
    iniFile := config.NewINIFile("config.ini")
```

Next, create a **Config** object to manage the application's providers. 
Since this application only has one provider, `iniFile`, pass a list with only one object into the constructor. 
```
    c := config.NewConfig([]config.Provider{iniFile})
```

It is a good idea to call the `Load()` function on your config object after creating it. 
```
    if err := c.Load(); err != nil{
        log.Fatal(err)
    }
```

However, calling the `Load()` function is not required. 
Each time a setting is requested, the `Load()` function is called first. 
If you are concerned about performance of calling `Load()` so frequently, see the [Advanced](#advanced) section.  

The `Config` object can be used to lookup settings from your providers. 
When performing a lookup on a setting, the key is a period-delimited string. 
For `ini` files, this means lookups are performed with `<section>.<item>` keys. 
For example, the settings in `config.json` can be looked up with the following keys:
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
All of the settings in your application are stored as strings. 
As shown above, the `Config` object has convenience functions for type conversions: `Int(), Float(), Bool()`. 

## Validation
The `Config` object has an optional `Validate` function that is called after the settings are loaded from the providers. 
Taking the `config.ini` example from above, create a `Validate` function that makes sure the `global.timeout` setting exists: 
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
Building off of the [Single Provider Example](#single-provider-example) above, this example will add environment variables as a means of configuration in addition to the `config.ini` file.
If they are set, the environment variables will override settings in `config.go`. 

As before, create a provider that will be responsible for loading configuration values from the `config.ini` file:
```
    iniFile := config.NewINIFile("config.ini")
```

Next, create an `Environment` provider that will be responsible for loading configuration values from environment variables. 
The `Environment` provider takes a map that associates setting keys with environment variables.
Since the environment variables should override the same settings keys as `config.ini`, construct the map like so:
```
    mappings := map[string]string{
        "APP_TIMEOUT": "global.timeout",
        "APP_FREQUENCY": "global.frequency",
        "APP_TIMEZONE": "local.time_zone",
        "APP_ENABLED": "local.enabled",
    }
    
    env := config.NewEnvironment(mappings)
```
Since there are two providers, add both of them to the `Config` object. 
The position of the provider in the list determines the ordering of settings lookups. 
Since environment variables should override the values in `config.ini`, put the `Environment` provider later in the list:
```
    providers := []config.Providers{iniFile, env}
    c := config.NewConfig(providers)
```

## CLI Provider Example
In addition to files and environment variables, applications tend to use command line arguments for configuration.
One of the most popular command line tool for golang is [Jeremy Saenz's CLI](https://github.com/codegangsta/cli).
This tool does an excellent job of allowing users to configure their applications using 
[flags](https://github.com/codegangsta/cli#flags) and [environment variables](https://github.com/codegangsta/cli#values-from-the-environment).
This works great for many applications, but can easily become messy when settings need to be loaded from other sources.
The `CLI` provider aims to make configuration management from any number of providers as simple as possible.


For the following example, assume the application includes a configuration file, `config.yaml`, and `main.go` with the following content:

### Config.yaml
```
message: "Hello from config.yaml"
silent: false
```

### Main.go
```
package main 

import (
    "github.com/codegangsta/cli"
    "github.com/zpatrick/go-config"
    "log"
    "os"
)

func initConfig() *config.Config {
    yamlFile := config.NewYAMLFile("config.yaml")
    return config.NewConfig([]config.Provider{yamlFile})
}

func main() {
    conf := initConfig()

    app := cli.NewApp()
    app.Flags = []cli.Flag{
        cli.StringFlag{
            Name:  "message",
            Value: "Hello from main.go",
            Usage: "Message to print",
        },
        cli.BoolFlag{
            Name: "silent",
            Usage: "Don't print the message",
        },
    }

    app.Action = func(c *cli.Context) {
        conf.Providers = append(conf.Providers, config.NewCLI(c, false))

        message, err := conf.String("message")
        if err != nil {
            log.Fatal(err)
        }

        silent, err := conf.Bool("silent")
        if err != nil {
            log.Fatal(err)
        }

        if !silent {
            log.Println(message)
        }
    }

    app.Run(os.Args)
}
```

### Creating the YAML Provider
The following lines create the `YAML` provider and the `Config` object
```
yamlFile := config.NewYAMLFile("config.yaml")
return config.NewConfig([]config.Provider{yamlFile})
```
Since this application uses `config.yaml` for configuration, the `YAMLFile` is used to load settings from that file.
This will load the settings `message="Hello from config.yaml"` and `silent=false`.

### Creating the CLI Provider
In the `app.Action` function, a new `CLI` provider is created to load settings. The `CLI` provider takes a `*cli.Context` and boolean argument `useDefaults`. Having default values for flags is useful, but unlike other flags, boolean flags always have a default value. Since this could lead to unwanted or unexpected behavior, users must specify which setting to use:
* If `useDefaults=true`, flags with default values will be loaded.
In context of this example, the `CLI` provider will load the settings `message="Hello from main.go"` and `silent=false` when the user runs the application without any arguments. 
These settings would overwrite the `message` and `silent` settings loaded from `config.yaml`.
* If `useDefaults=false`, only flags that have been set via the command line will be loaded as settings.
In context of this example, the `CLI` provider will not load any settings when the user runs the application without any arguments.
This allows the setting loaded by `config.yaml` to not be overwitten.

The following line creates the `CLI` provider and appends it to the existing providers:
```
conf.Providers = append(conf.Providers, config.NewCLI(c, false))
```
Note that `useDefaults=false`. This way, settings in `config.yaml` aren't overwitten by the default flag values.

### Running the Application
**No Arguments**: The `message` and `silent` settings from `config.yaml` are used
```
> go run main.go
Hello from config.yaml
```
**Message Flag**: The `message` setting from the `CLI` provider is used. The `silent` setting from `config.yaml` is used
```
> go run main.go --message "Hello from the command line"
Hello from the command line
```
**Silent and Message Flag**: The `message` and `silent` settings from the `CLI` provider are used
```
> go run main.go --message "this shouldn't print" --silent
<no output>
```

**No Arguments with `useDefaults=True`**: To further demonstrate the different behavior of `useDefault`, here is the output when `useDefault=true`. The default `message` and `silent` settings from the `CLI` are used:
```
> go run main.go
Hello from main.go
```

# Advanced

## Defaults
Managing default values for settings can be accomplished multiple ways:
* A configuration file that contains all of the defaults
* Using `Or` functions with defaults (e.g. `StringOr(...)`, `IntOr(...)`, etc.)
* Using the `Static` provider

The `Static` provider takes key value mappings for settings and simply returns those values when `Load()` is called. 
This is a nice pattern as it doesn't require additional configuration files and it places all defaults into a single place in your code. 
Make sure to set your defaults as the first provider in your application so they can be overridden by other providers:
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
Each time a lookup is performed, the `Load()` function is called on each provider. 
This can lead to poor performance and be unecessary for certain providers. 
There are two built in objects which change how frequently loads are performed:
* `OnceLoader` - Loads the provider's settings one time
* `CachedLoader` - Loads the provider's settings at least one time and caches the results. 
The `Invalidate()` function can be called to force a new load next time a lookup is performed.

Building off of the [Multiple Providers Example](#multiple-providers-example) above, 
use the `OnceLoader` for the `iniFile` provider and keep the default behavior for the `Environment` provider (perform a load each time a lookup is requested):
```
    env := config.NewEnvironment(...)
    iniFile := config.NewINIFile("config.ini")
    iniFileLoader := config.NewOnceLoader(iniFile)
    
    providers := []config.Provider{iniFileLoader, env}
    c := config.NewConfig(providers)
```
The first time a lookup is performed, the provider's `Load()` function will be called. 
All other calls will use the same settings as the original lookup.  

The `CachedLoader` behaves in a similar manner except that it contains an `Invalidate()` function. 
After `Invalidate()` is called, the provider's `Load()` function will be executed the next time a lookup is performed.
```
    env := config.NewEnvironment(...)
    iniFile := config.NewINIFile("config.ini")
    iniFileLoader := config.NewCachedLoader(iniFile)
    
    providers := []config.Provider{iniFileLoader, env}
    c := config.NewConfig(providers)

    ...
    // will execute iniFile.Load() next time a lookup is performed
    iniFileLoader.Invalidate()
    ...
```

## Resolvers
Sometimes, keeping setting keys consistent between different files isn't possible. 
For example, say the file `config.json` contained:
```
{
    "items": {
        "server": {
            "timeout": 30
        }
    }
}
```
And another file `config.ini` contained:
```
[server]
timeout=30
```

The `timeout` setting would have the key `items.server.timeout` for the `json` file and `server.timeout` for the `ini` file when they are actually intended to reference the same setting. 
A `Resolver` can be used to change the mappings in a provider. 
For example, wrap the `json` provider in a `Resolver` in order to resolve `items.server.timeout` as `server.timeout`:
```
    iniFile := config.NewINIFile("config.ini")
    jsonFile := config.NewJSONFile("config.json")
    mappings := map[string]string{
        "items.server.timeout": "server.timeout",
    }

    JSONFileResolver := config.NewResolver(jsonFile, mappings)

    providers := []config.Provider{iniFile, jsonFileResolver}
    c := config.NewConfig(providers)
```

The canonical key for the setting is now `server.timeout`

## Custom Providers
Custom providers must fulfill the `Provider` interface:
```
    type Provider interface {
        Load() (map[string]string, error)
    }
```

The `Load()` function returns settings as key/value pairs. 
Providers flatten namespaces using period-delimited strings. 
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
This is allows providers to override and lookup settings using a canonical key. 


# License
This work is published under the MIT license.

Please see the `LICENSE` file for details.