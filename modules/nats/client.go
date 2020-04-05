package nats

import (
	"github.com/nats-io/stan.go"
	"nats_webhook/modules/config"
)

// Connect ...
func Connect() (stan.Conn, error) {
	return stan.Connect(config.Vars.NATS.ClusterID, config.Vars.NATS.ClientID, stan.NatsURL(config.Vars.NATS.Endpoint))
}
