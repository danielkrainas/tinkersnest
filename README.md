# TinkersNest

[![License](https://img.shields.io/badge/license-Unlicense-blue.svg?style=flat)](UNLICENSE) [![Go Report Card](https://goreportcard.com/badge/github.com/danielkrainas/tinkersnest)](https://goreportcard.com/report/github.com/danielkrainas/tinkersnest) [![Docker Hub](https://img.shields.io/docker/pulls/dakr/tinkersnest.svg?style=flat)](https://hub.docker.com/r/dakr/tinkersnest/)

TinkersNest is a developer-centric backend CMS service. 

## Installation

> $ go get github.com/danielkrainas/tinkersnest

## Usage

> $ tinkersnest [command] <config_path>

Most commands require a configuration path provided as an argument or in the `TINKERS_CONFIG_PATH` environment variable. 

### Serve mode

This is the primary mode for TinkersNest. It hosts the HTTP API server.

> $ tinkersnest serve <config_path>

**Example** - with the default config:

> $ tinkersnest serve ./config.default.yml

## Configuration

A configuration file is *required* for TinkersNest but environment variables can be used to override configuration. A configuration file can be specified as a parameter or with the `TINKERS_CONFIG_PATH` environment variable. 

All configuration environment variables are prefixed by `TINKERS_` and the paths are separated by an underscore(`_`). Some examples:

- `TINKERS_LOGGING_LEVEL=warn`
- `TINKERS_HTTP_ADDR=localhost:2345`
- `TINKERS_STORAGE_INMEMORY=true`

A development configuration file is included: `/config.dev.yml` and a `/config.local.yml` has already been added to gitignore to be used for local testing or development.

```yaml
# configuration schema version number, only `0.1`
version: 0.1

# log stuff
logging:
  # minimum event level to log: `error`, `warn`, `info`, or `debug`
  level: 'debug'
  # log output format: `text` or `json`
  formatter: 'text'
  # custom fields to be added and displayed in the log
  fields:
    customfield1: 'value'

# http server stuff
http:
  # host:port address for the server to listen on
  addr: ':9240'
  # http host
  host: 'localhost'

  # CORS stuff
  cors:
    # origins to allow
    origins: ['http://localhost:5555']
    # methods to allow
    methods: ['GET','POST','OPTIONS','DELETE','CONNECT']
    # headers to allow
    headers: ['*']

# storage driver and parameters
storage:
  inmemory:
    param1: 'val'

# the in-memory driver has no parameters so it can be declared as a string
storage: 'inmemory'
```

`storage` only allows specification of *one* driver per configuration. Any additional ones will cause a validation error when the application starts.

## Bugs and Feedback

If you see a bug or have a suggestion, feel free to open an issue [here](https://github.com/danielkrainas/tinkersnest/issues).

## Contributions

PR's welcome! There are no strict style guidelines, just follow best practices and try to keep with the general look & feel of the code present. All submissions should atleast be `go fmt -s` and have a test to verify *(if applicable)*.

For details on how to extend and develop TinkersNest, see the [dev documentation](docs/development/).

## License

[Unlicense](http://unlicense.org/UNLICENSE). This is a Public Domain work. 

[![Public Domain](https://licensebuttons.net/p/mark/1.0/88x31.png)](http://questioncopyright.org/promise)

> ["Make art not law"](http://questioncopyright.org/make_art_not_law_interview) -Nina Paley
