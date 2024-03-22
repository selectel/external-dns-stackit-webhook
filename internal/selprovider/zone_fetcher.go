package selprovider

import (
	"context"
	"strconv"

	domains "github.com/selectel/domains-go/pkg/v2"
	"sigs.k8s.io/external-dns/endpoint"
)

type zoneFetcher struct {
	domainFilter endpoint.DomainFilter
}

func newZoneFetcher(
	domainFilter endpoint.DomainFilter,
) *zoneFetcher {
	return &zoneFetcher{
		domainFilter: domainFilter,
	}
}

// zones returns filtered list of v2.Zone if domainFilter is set.
func (z *zoneFetcher) zones(ctx context.Context, client domains.DNSClient[domains.Zone, domains.RRSet]) ([]*domains.Zone, error) {
	if len(z.domainFilter.Filters) == 0 {
		zones, err := z.fetchZones(ctx, client, map[string]string{})
		if err != nil {
			return nil, err
		}

		return zones, nil
	}

	var result []*domains.Zone
	// send one request per filter
	for _, filter := range z.domainFilter.Filters {
		zones, err := z.fetchZones(ctx, client, map[string]string{
			"filter": filter,
		})
		if err != nil {
			return nil, err
		}
		result = append(result, zones...)
	}

	return result, nil
}

// fetchZones fetches all []v2.Zone from Selectel DNS API. It may be filtered with options["filter"] provided.
func (z *zoneFetcher) fetchZones(
	ctx context.Context,
	client domains.DNSClient[domains.Zone, domains.RRSet],
	options map[string]string,
) ([]*domains.Zone, error) {
	options["limit"] = "1000"
	options["offset"] = "0"

	var zones []*domains.Zone

	for {
		zonesResponse, err := client.ListZones(ctx, &options)
		if err != nil {
			return nil, err
		}

		zones = append(zones, zonesResponse.GetItems()...)

		options["offset"] = strconv.Itoa(zonesResponse.GetNextOffset())
		if zonesResponse.GetNextOffset() == 0 {
			break
		}
	}

	return zones, nil
}
