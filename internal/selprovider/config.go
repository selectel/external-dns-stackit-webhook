package selprovider

import (
	"sigs.k8s.io/external-dns/endpoint"
)

// Config is used to configure the creation of the Provider.
type Config struct {
	BaseURL          string
	KeystoneProvider KeystoneProvider
	DomainFilter     endpoint.DomainFilter
	DryRun           bool
	Workers          int
}

//go:generate mockgen -destination=./mock/keystone_provider.go -source=./config.go KeystoneProvider
type KeystoneProvider interface {
	GetToken() (string, error)
}
