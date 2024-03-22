package selprovider

import (
	"context"
	"fmt"
	"strconv"

	domains "github.com/selectel/domains-go/pkg/v2"
	"go.uber.org/zap"
	"sigs.k8s.io/external-dns/endpoint"
)

type rrSetFetcher struct {
	domainFilter endpoint.DomainFilter
	logger       *zap.Logger
}

func newRRSetFetcher(
	domainFilter endpoint.DomainFilter,
	logger *zap.Logger,
) *rrSetFetcher {
	return &rrSetFetcher{
		domainFilter: domainFilter,
		logger:       logger,
	}
}

// fetchRecords fetches all []v2.RRSet from Selectel DNS API for given zone id in options["filter"].
func (r *rrSetFetcher) fetchRecords(
	ctx context.Context,
	client domains.DNSClient[domains.Zone, domains.RRSet],
	zoneId string,
	options map[string]string,
) ([]*domains.RRSet, error) {
	options["limit"] = "1000"
	options["offset"] = "0"

	var rrSets []*domains.RRSet

	for {
		rrSetsResponse, err := client.ListRRSets(ctx, zoneId, &options)
		if err != nil {
			return nil, err
		}

		rrSets = append(rrSets, rrSetsResponse.GetItems()...)

		options["offset"] = strconv.Itoa(rrSetsResponse.GetNextOffset())
		if rrSetsResponse.GetNextOffset() == 0 {
			break
		}
	}

	return rrSets, nil
}

// getRRSetForUpdateDeletion returns the record set to be deleted and the zone it belongs to.
func (r *rrSetFetcher) getRRSetForUpdateDeletion(
	ctx context.Context,
	client domains.DNSClient[domains.Zone, domains.RRSet],
	change *endpoint.Endpoint,
	zones []*domains.Zone,
) (*domains.Zone, *domains.RRSet, error) {
	resultZone, found := findBestMatchingZone(change.DNSName, zones)
	if !found {
		r.logger.Error(
			"record set name contains no zone dns name",
			zap.String("name", change.DNSName),
		)

		return nil, nil, fmt.Errorf("record set name contains no zone dns name")
	}

	domainRrSets, err := r.fetchRecords(ctx, client, resultZone.ID, map[string]string{
		"name": change.DNSName,
	})
	if err != nil {
		return nil, nil, err
	}

	resultRRSet, found := findRRSet(change.DNSName, change.RecordType, domainRrSets)
	if !found {
		r.logger.Info("record not found on record sets", zap.String("name", change.DNSName))

		return nil, nil, fmt.Errorf("record not found on record sets")
	}

	return resultZone, resultRRSet, nil
}
