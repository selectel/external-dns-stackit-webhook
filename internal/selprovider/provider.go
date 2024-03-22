package selprovider

import (
	"net/http"

	domains "github.com/selectel/domains-go/pkg/v2"
	"github.com/selectel/external-dns-stackit-webhook/pkg/httpclient"
	"go.uber.org/zap"
	"sigs.k8s.io/external-dns/endpoint"
	"sigs.k8s.io/external-dns/provider"
)

// Provider implements the DNS provider interface for Selectel DNS.
type Provider struct {
	provider.BaseProvider
	domainFilter       endpoint.DomainFilter
	keystoneProvider   KeystoneProvider
	endpoint           string
	dryRun             bool
	workers            int
	logger             *zap.Logger
	zoneFetcherClient  *zoneFetcher
	rrSetFetcherClient *rrSetFetcher
}

// getDomainsClient returns v2.DNSClient with provided keystone and user-agent from httpclient.DefaultUserAgent.
func (p *Provider) getDomainsClient() (domains.DNSClient[domains.Zone, domains.RRSet], error) {
	token, err := p.keystoneProvider.GetToken()
	if err != nil {
		p.logger.Error("authorization error during getting keystone token", zap.Error(err))

		return nil, err
	}

	httpClient := httpclient.Default()
	headers := http.Header{}
	headers.Add("X-Auth-Token", token)
	headers.Add("User-Agent", httpclient.DefaultUserAgent)

	return domains.NewClient(p.endpoint, &httpClient, headers), nil
}

// New creates a new Selectel DNS provider.
func New(config Config, logger *zap.Logger) (*Provider, error) {
	return &Provider{
		domainFilter:       config.DomainFilter,
		dryRun:             config.DryRun,
		workers:            config.Workers,
		logger:             logger,
		keystoneProvider:   config.KeystoneProvider,
		endpoint:           config.BaseURL,
		zoneFetcherClient:  newZoneFetcher(config.DomainFilter),
		rrSetFetcherClient: newRRSetFetcher(config.DomainFilter, logger),
	}, nil
}
