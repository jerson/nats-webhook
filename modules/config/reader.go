// Package config ...
package config

import (
	"github.com/jinzhu/configor"
	"os"
	"path/filepath"
)

// Server ...
type Server struct {
	Port int `toml:"port" default:"8080"`
}

// API ...
type API struct {
	Key string `toml:"key" default:""`
}

// NATS ...
type NATS struct {
	Endpoint  string `toml:"endpoint" default:"nats://localhost:4222"`
	ClientID  string `toml:"client_id" default:"webhook"`
	ClusterID string `toml:"cluster_id" required:"true"`
}

// Vars ...
var Vars = struct {
	Debug  bool
	API    API    `toml:"api"`
	Server Server `toml:"server"`
	NATS   NATS   `toml:"nats"`
}{}

// ReadDefault ...
func ReadDefault() error {
	file, err := filepath.Abs("./config.toml")
	if err != nil {
		return err
	}
	return Read(file)
}

// Read ...
func Read(file string) error {

	debug := os.Getenv("DEBUG") == "true"
	err := configor.New(&configor.Config{
		ENVPrefix: "APP",
		Debug:     debug,
		Verbose:   false,
	}).Load(&Vars, file)
	if err != nil {
		return err
	}

	Vars.Debug = debug
	return nil
}
