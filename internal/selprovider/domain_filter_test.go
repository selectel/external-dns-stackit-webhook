package selprovider

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/external-dns/endpoint"
)

func TestGetDomainFilter(t *testing.T) {
	t.Parallel()

	server := getServerRecords(t)
	defer server.Close()

	dnsProvider, err := getDefaultTestProvider(server, getDefaultKeystoneProvider(t, 0))
	assert.NoError(t, err)

	domainFilter := dnsProvider.GetDomainFilter()
	assert.Equal(t, domainFilter, endpoint.DomainFilter{})
}
