# NATS Webhook

Simple Webhook server and store in `NATS streaming` for use later.

## Usage

Just use this endpoint for a Webhook

```bash
http://localhost:8080/event/{source}/{subject}
```

- `{source}`: just a name like "mysource" its for identify later on `NATS`
- `{subject}`: subject event for Publish to `NATS`

Later you can subscribe to any subject a usual its on `NATS streaming`

```go
package main
import (
	"fmt"
  "time"
	"github.com/nats-io/stan.go"
)

func main(){
    sc, err := stan.Connect("default", "itsme", stan.NatsURL("nats://localhost:4222"))
    if err != nil {
      panic(err)
    }
    defer sc.Close()

    sub, err := sc.Subscribe("mysubject", func(m *stan.Msg) {
      fmt.Printf("Hi: %s\n", string(m.Data))
      m.Ack()
    })
    if err != nil {
      panic(err)
    }
    defer sub.Unsubscribe()

    // delay for wait some seconds... its not really required
    time.Sleep(time.Second * 10)
}
```

**IMPORTANT**: `m.Data` is json.Marshal of `Payload` struct

```go
type Payload struct {
	ID      string `json:"id"`
	Source  string `json:"source"`
	Subject string `json:"subject"`
	Body    []byte `json:"body"`
}
```

## Deploy

### Config

Before deploy you should know all vars, you have two options:

#### config.toml

```toml
[server]
port=8080

[api]
key="sampletoken"

[nats]
endpoint="nats://localhost:4222"
client_id="webhook"
cluster_id="nats"
```

#### Environments vars

```bash
export APP_SERVER_PORT=8080
export APP_API_KEY=sampletoken
export APP_NATS_ENDPOINT=nats://localhost:4222
export APP_NATS_CLIENTID=webhook
export APP_NATS_CLUSTERID=nats
```

### Authorization

if you define `APP_API_KEY` env or `key="sampletoken"` in `config.toml` you should send `Authorization` header like this `Authorization: Bearer {sampletoken}`

### Docker

You can use a `docker-compose.yml` like this

```yaml
version: "3"
services:
  webhook:
    image: registry.gitlab.com/pardacho/nats-webhook:latest
    environment:
      - DEBUG=true
      - APP_NATS_ENDPOINT=nats://nats:4222
      - APP_NATS_CLUSTERID=default
      - APP_API_KEY=sampletoken
    ports:
      - 8080:8080
  nats:
    image: nats-streaming:latest
    volumes:
      - storage:/data
    command:
      - "--cluster_id"
      - "default"
      - "-store"
      - "file"
      - "-dir"
      - "/data"
      - "-p"
      - "4222"
    ports:
      - 4222:4222
volumes:
  storage:
```

### Binaries

just go to `Releases` section, download an run in your server

<https://github.com/jerson/nats-webhook/releases>

### Script

```bash
curl -sfL https://raw.githubusercontent.com/jerson/nats-webhook/master/install.sh | sh
```
