package selprovider

import (
	"reflect"
	"testing"

	domains "github.com/selectel/domains-go/pkg/v2"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"sigs.k8s.io/external-dns/endpoint"
)

//nolint:funlen
func TestModifyChange(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		ep   *endpoint.Endpoint
		want *endpoint.Endpoint
	}{
		{
			name: "trailing dot in DNSName",
			ep: &endpoint.Endpoint{
				DNSName:    "example.com",
				RecordType: "A",
				Targets:    []string{"1.2.3.4"},
				RecordTTL:  endpoint.TTL(300),
			},
			want: &endpoint.Endpoint{
				DNSName:    "example.com.",
				RecordType: "A",
				Targets:    []string{"1.2.3.4"},
				RecordTTL:  endpoint.TTL(300),
			},
		},
		{
			name: "ttl added",
			ep: &endpoint.Endpoint{
				DNSName:    "example.com.",
				RecordType: "A",
				Targets:    []string{"1.2.3.4"},
			},
			want: &endpoint.Endpoint{
				DNSName:    "example.com.",
				RecordType: "A",
				Targets:    []string{"1.2.3.4"},
				RecordTTL:  endpoint.TTL(300),
			},
		},
		{
			name: "trailing dot in CNAME targets",
			ep: &endpoint.Endpoint{
				DNSName:    "example.com.",
				RecordType: "CNAME",
				Targets:    []string{"sub.example.com"},
				RecordTTL:  endpoint.TTL(300),
			},
			want: &endpoint.Endpoint{
				DNSName:    "example.com.",
				RecordType: "CNAME",
				Targets:    []string{"sub.example.com."},
				RecordTTL:  endpoint.TTL(300),
			},
		},
		{
			name: "trailing dot in ALIAS targets",
			ep: &endpoint.Endpoint{
				DNSName:    "example.com.",
				RecordType: "ALIAS",
				Targets:    []string{"sub.example.com"},
				RecordTTL:  endpoint.TTL(300),
			},
			want: &endpoint.Endpoint{
				DNSName:    "example.com.",
				RecordType: "ALIAS",
				Targets:    []string{"sub.example.com."},
				RecordTTL:  endpoint.TTL(300),
			},
		},
		{
			name: "trailing dot in MX targets",
			ep: &endpoint.Endpoint{
				DNSName:    "example.com.",
				RecordType: "MX",
				Targets:    []string{"mail.example.com"},
				RecordTTL:  endpoint.TTL(300),
			},
			want: &endpoint.Endpoint{
				DNSName:    "example.com.",
				RecordType: "MX",
				Targets:    []string{"mail.example.com."},
				RecordTTL:  endpoint.TTL(300),
			},
		},
		{
			name: "trailing dot in SRV targets",
			ep: &endpoint.Endpoint{
				DNSName:    "_xmpp._tcp.example.com.",
				RecordType: "SRV",
				Targets:    []string{"sub.example.com"},
				RecordTTL:  endpoint.TTL(300),
			},
			want: &endpoint.Endpoint{
				DNSName:    "_xmpp._tcp.example.com.",
				RecordType: "SRV",
				Targets:    []string{"sub.example.com."},
				RecordTTL:  endpoint.TTL(300),
			},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			modifyChange(tc.ep)
			assert.EqualValues(t, tc.want, tc.ep)
		})
	}
}

func TestGetRRSetRecordPost(t *testing.T) {
	t.Parallel()

	change := &endpoint.Endpoint{
		DNSName:    "test.",
		RecordTTL:  endpoint.TTL(300),
		RecordType: "A",
		Targets: endpoint.Targets{
			"192.0.2.1",
			"192.0.2.2",
		},
	}
	expected := &domains.RRSet{
		Name: "test.",
		TTL:  300,
		Type: "A",
		Records: []domains.RecordItem{
			{
				Content: "192.0.2.1",
			},
			{
				Content: "192.0.2.2",
			},
		},
	}
	got := getRRSetRecord(change)
	if !reflect.DeepEqual(got, expected) {
		t.Errorf("getRRSetRecord() = %v, want %v", got, expected)
	}
}

func TestFindBestMatchingZone(t *testing.T) {
	t.Parallel()

	zones := []*domains.Zone{
		{Name: "foo.com"},
		{Name: "bar.com"},
		{Name: "baz.com"},
	}

	tests := []struct {
		name      string
		rrSetName string
		want      *domains.Zone
		wantFound bool
	}{
		{"Matching Zone", "www.foo.com", zones[0], true},
		{"No Matching Zone", "www.test.com", nil, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, found := findBestMatchingZone(tt.rrSetName, zones)
			if !reflect.DeepEqual(got, tt.want) || found != tt.wantFound {
				t.Errorf("findBestMatchingZone() = %v, %v, want %v, %v", got, found, tt.want, tt.wantFound)
			}
		})
	}
}

func TestFindRRSet(t *testing.T) {
	t.Parallel()

	rrSets := []*domains.RRSet{
		{Name: "www.foo.com", Type: "A"},
		{Name: "www.bar.com", Type: "A"},
		{Name: "www.baz.com", Type: "A"},
	}

	tests := []struct {
		name       string
		rrSetName  string
		recordType string
		want       *domains.RRSet
		wantFound  bool
	}{
		{"Matching RRSet", "www.foo.com", "A", rrSets[0], true},
		{"No Matching RRSet", "www.test.com", "A", nil, false},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, found := findRRSet(tt.rrSetName, tt.recordType, rrSets)
			if !reflect.DeepEqual(got, tt.want) || found != tt.wantFound {
				t.Errorf("findRRSet() = %v, %v, want %v, %v", got, found, tt.want, tt.wantFound)
			}
		})
	}
}

func TestGetLogFields(t *testing.T) {
	t.Parallel()

	change := &endpoint.Endpoint{
		DNSName:    "test.",
		RecordTTL:  endpoint.TTL(300),
		RecordType: "A",
		Targets: endpoint.Targets{
			"192.0.2.1",
			"192.0.2.2",
		},
	}

	expected := []zap.Field{
		zap.String("record", "test."),
		zap.String("content", "192.0.2.1,192.0.2.2"),
		zap.String("type", "A"),
		zap.String("action", "create"),
		zap.String("id", "123"),
	}

	got := getLogFields(change, "create", "123")

	if !reflect.DeepEqual(got, expected) {
		t.Errorf("getLogFields() = %v, want %v", got, expected)
	}
}
