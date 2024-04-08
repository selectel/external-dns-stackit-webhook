package selprovider

import (
	"strings"

	domains "github.com/selectel/domains-go/pkg/v2"
	"go.uber.org/zap"
	"sigs.k8s.io/external-dns/endpoint"
)

// findBestMatchingZone finds the best matching zone for a given record set name. The criteria are
// that the zone name is contained in the record set name and that the zone name is the longest
// possible match. Eg foo.bar.com. would have prejudice over bar.com. if rr set name is foo.bar.com.
func findBestMatchingZone(rrSetName string, zones []*domains.Zone) (*domains.Zone, bool) {
	count := 0
	var domainZone *domains.Zone
	for _, zone := range zones {
		if len(zone.Name) > count && strings.Contains(appendDotIfNotExists(rrSetName), zone.Name) {
			count = len(zone.Name)
			domainZone = zone
		}
	}

	if count == 0 {
		return nil, false
	}

	return domainZone, true
}

// findRRSet finds a record set by name and type in a list of record sets.
func findRRSet(rrSetName, recordType string, rrSets []*domains.RRSet) (*domains.RRSet, bool) {
	for _, rrSet := range rrSets {
		if rrSet.Name == rrSetName && string(rrSet.Type) == recordType {
			return rrSet, true
		}
	}

	return nil, false
}

// appendDotIfNotExists appends a dot to the end of a string if it doesn't already end with a dot.
func appendDotIfNotExists(s string) string {
	if !strings.HasSuffix(s, ".") {
		return s + "."
	}

	return s
}

// modifyChange modifies a change to ensure it is valid for this provider.
func modifyChange(change *endpoint.Endpoint) {
	change.DNSName = appendDotIfNotExists(change.DNSName)

	if change.RecordTTL == 0 {
		change.RecordTTL = 300
	}
}

// getRRSetRecord returns a v2.RRSet from a change for the api client.
func getRRSetRecord(change *endpoint.Endpoint) *domains.RRSet {
	records := make([]domains.RecordItem, len(change.Targets))
	for i, target := range change.Targets {
		records[i] = domains.RecordItem{
			Content: target,
		}
	}

	return &domains.RRSet{
		Name:    change.DNSName,
		Records: records,
		TTL:     int(change.RecordTTL),
		Type:    domains.RecordType(change.RecordType),
	}
}

// getLogFields returns a log.Fields object for a change.
func getLogFields(change *endpoint.Endpoint, action string, id string) []zap.Field {
	return []zap.Field{
		zap.String("record", change.DNSName),
		zap.String("content", strings.Join(change.Targets, ",")),
		zap.String("type", change.RecordType),
		zap.String("action", action),
		zap.String("id", id),
	}
}
