# Comic
![GitHub Workflow Status](https://github.com/zaininfo/comic/workflows/CI/badge.svg)

A library to manage application configurations, especially useful for multi-command applications, where each command needs a different subset of configurations.

It leverages [Viper](https://github.com/spf13/viper) for handling configurations, while adding support for:
- marking required configurations within configuration files
- using different subsets of a common configuration for various commands of an application [e.g.](#multi-command-applications)

## Usage

### Example
**Configuration file (`config.yaml`):**
```yaml
# all configurations (with default values, if any)
server:
  host: localhost
  port:
ttl: 30s

# names of required configurations (for a single command application)
required:
  main:
    server:
      host:
      port:
```

**Environment:**
```sh
export SERVER_PORT=80
```

**Code:**
```go
package main

import (
	"time"

	"github.com/zaininfo/comic"
)

type Config struct {
	Server Server        `mapstructure:"SERVER"`
	TTL    time.Duration `mapstructure:"TTL"`
}

type Server struct {
	Host string `mapstructure:"HOST"`
	Port int    `mapstructure:"PORT"`
}

func mustLoad() Config {
	var cfg Config
	comic.MustLoad(&cfg)

	return cfg
}

func main() {
	// cfg := mustLoad()
}
```

Please, check Viper documentation for all supported [config file types](https://github.com/spf13/viper#reading-config-files) and [config data types](https://github.com/spf13/viper#getting-values-from-viper).

Note that [`mapstructure`](https://github.com/mitchellh/mapstructure) tags are used to unmarshal configuration data.

### Options
The following options can be used to change the behavior of Comic.

| Name                     | Default               | Description                                                                                                     |
|--------------------------|:---------------------:|-----------------------------------------------------------------------------------------------------------------|
| ConfigFileName           | config                | The name of the configuration file (without extension, but actual file name should have appropriate extension). |
| ConfigFilePath           | . (working directory) | The path to the configuration file.                                                                             |
| SingleCommandAppName     | main                  | The name used in the `required` section of the configuration file for a single command application.             |
| EnvVarNestedKeySeparator | _                     | The separator used for referring to nested environment variables.                                               |

### Functions
- `New()`
  - It returns a new instance of Comic with default options.
- `NewWithOptions(opts Options)`
  - It returns a new instance of Comic with supplied options.
- `FromCommandPath(commandPath string)`
  - It removes the binary name from the supplied command path and returns the rest of it.
- `Viper()`
  - It returns the Viper instance in use by Comic, which is unique for package-level exported Comic and all instances of Comic.
- `MustLoad(cfg interface{})`
  - It loads configurations from file & environment into `cfg` after verifying all required configurations; it panics on failure.
- `MustLoadForCommand(cfg interface{}, commandName string)`
  - It loads configurations from file & environment into `cfg` after verifying all required configurations of `commandName`; it panics on failure.
- `Load(cfg interface{})`
  - Same as `MustLoad(cfg interface{})`, but returns an error on failure.
- `LoadForCommand(cfg interface{}, commandName string)`
  - Same as `MustLoadForCommand(cfg interface{}, commandName string)`, but returns an error on failure.

The `Viper()` & all `*Load*()` functions can be called on both package-level exported Comic and an instance of Comic.

**Important:** the configuration structure passed to any of the `*Load*()` functions should be a pointer.

## Multi-command applications

### Example
**Configuration file (`config.yaml`):**
```yaml
# all configurations (with default values, if any)
server:
  host: localhost
  port:
ttl: 30s
db_connections:

# names of required configurations (for all commands of the application)
required:
  api:
    server:
      host:
      port:
    db_connections:
  indexer:
    db_connections:
```

**Environment (`api`):**
```sh
export SERVER_PORT=80
export DB_CONNECTIONS=5
```

**Environment (`indexer`):**
```sh
export DB_CONNECTIONS=5
```

**Code:**
```go
package config

import (
	"time"

	"github.com/zaininfo/comic"
)

type Config struct {
	Server        Server        `mapstructure:"SERVER"`
	TTL           time.Duration `mapstructure:"TTL"`
	DBConnections int           `mapstructure:"DB_CONNECTIONS"`
}

type Server struct {
	Host string `mapstructure:"HOST"`
	Port int    `mapstructure:"PORT"`
}

func mustLoad(commandName string) Config {
	var cfg Config
	comic.MustLoadForCommand(&cfg, commandName)

	return cfg
}

func main() {
	// apiCfg := mustLoad("api")
	// indexerCfg := mustLoad("indexer")
}
```

## Q&A

Q: What's with it being comical?

A: It's **C**~~onfiguration~~ ~~F~~**o**~~r~~ **M**~~ult~~**i**~~ple~~ **C**~~ommand~~ ~~Applications~~.
