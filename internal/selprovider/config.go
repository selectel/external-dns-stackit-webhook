package selprovider

import (
	"sigs.k8s.io/external-dns/endpoint"
)

// Config is used to configure the creation of the Provider.
type Config struct {
	// BaseURL is a Selectel DNS API endpoint for v2.DNSClient
	BaseURL string
	// KeystoneProvider needed to generate X-Auth-Token header with keystone-header for requests to the DNS API.
	KeystoneProvider KeystoneProvider
	// DomainFilter is a list with domains that will be affected. If it is empty all available domains will be affected.
	DomainFilter endpoint.DomainFilter
	// DryRun is a flag specifies user's wish to run without requests to the DNS API
	DryRun bool
	// Workers is a number of goroutines that will create requests to the DNS API.
	Workers int
}

//go:generate mockgen -destination=./mock/keystone_provider.go -source=./config.go KeystoneProvider
type KeystoneProvider interface {
	GetToken() (string, error)
}
