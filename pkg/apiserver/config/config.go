package config

import "github.com/1ch0/go-restful/pkg/apiserver/infrastructure/datastore"

// Config config for server
type Config struct {
	// api server bind address
	BindAddr string

	// Datastore config
	Datastore datastore.Config
}
