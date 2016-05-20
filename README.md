# [Wit.ai](https://wit.ai/) Go client

Conversational app from Wit.ai [quick start](https://wit.ai/docs/quickstart).
Try the console [example](https://github.com/marselester/witty-ai/#console-example).

```
$ WITTY_TOKEN="<YOUR-TOKEN>" go run example.go
> what's the weather?
```

API is similar to [Python SDK](https://github.com/wit-ai/pywit). Basic actions
are defined as

```go
ai := witty.NewClient(token, nil)
ai.SayAct = say
ai.MergeAct = merge
```

custom actions are set in `ai.Actions` map

```go
ai.Actions["fetch-weather"] = fetchWeather
```

## Console example

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "log"

    "github.com/marselester/witty-ai"
)

func main() {
    token := os.Getenv("WITTY_TOKEN")

    ai := witty.NewClient(token, nil)
    ai.MergeAct = merge
    ai.Actions["fetch-weather"] = fetchWeather

    sessID := "my-session-id"
    ctx := witty.Context{}
    var err error

    fmt.Print("> ")
    input := bufio.NewScanner(os.Stdin)
    for input.Scan() {
        userMsg := input.Text()
        ctx, err = ai.RunActions(sessID, userMsg, ctx)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Print("> ")
    }
}

func merge(sessID string, ctx witty.Context, entities witty.Entities) witty.Context {
    // Retrieve the location entity and store it into a context field.
    if _, ok := entities["location"]; ok {
        ctx["loc"] = entities["location"][0]["value"]
    }
    return ctx
}

func fetchWeather(sessID string, ctx witty.Context) witty.Context {
    ctx["forecast"] = "sunny"
    return ctx
}
```
