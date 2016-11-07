# TinkersNest API Client

Client library for the TinkersNest API. 

Supported Endpoints:

- Blog Posts


## Installation

> $ go get github.com/danielkrainas/tinkersnest/api/client


## Usage

How to instantiate a new client:

```go
package main

import (
	"net/http"
	"github.com/danielkrainas/tinkersnest/api/client"
)

// http/https url of the tinkersnest service
const ENDPOINT = "http://localhost:9240"

func main() {
	// Create a new client
	c := client.New(ENDPOINT, http.DefaultClient)
}
```


## Example

A more detailed example can be found [here.](https://github.com/danielkrainas/tinkersnest/tree/master/api/client/example)

