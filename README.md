# TGO Plugin Go SDK

Go SDK for developing plugins for the TGO (The Great Open) customer service system.

## Installation

```bash
go get github.com/tgoai/tgo-plugin-go
```

## Quick Start

```go
package main

import (
    "github.com/tgoai/tgo-plugin-go"
)

type MyPlugin struct {
    tgo.BasePlugin
}

func (p *MyPlugin) Name() string    { return "hello-go" }
func (p *MyPlugin) Version() string { return "1.0.0" }

func (p *MyPlugin) Capabilities() []tgo.Capability {
    return []tgo.Capability{
        tgo.VisitorPanel("Hello Panel"),
    }
}

func (p *MyPlugin) OnVisitorPanelRender(ctx *tgo.RenderContext) interface{} {
    return tgo.NewKeyValue("Go SDK").Add("Status", "Working!")
}

func main() {
    tgo.Run(&MyPlugin{})
}
```

## Local Debugging

When running TGO via Docker Compose, the plugin socket is mounted to `./data/tgo-api/run/tgo.sock`. You can connect your local plugin to this path for debugging:

```go
func main() {
    tgo.Run(&MyPlugin{}, tgo.WithSocketPath("./data/tgo-api/run/tgo.sock"))
}
```

## Features

- **Type Safe**: Native Go structs for all protocols and UI templates.
- **Easy UI Building**: Chainable builders for KeyValue, Table, Form, and more.
- **Standard Protocol**: Full implementation of TGO's JSON-RPC 2.0 over Unix Socket.

## Documentation

Full documentation is available at [https://tgo.ai/docs/plugin/overview](https://tgo.ai/docs/plugin/overview).

## License

MIT

