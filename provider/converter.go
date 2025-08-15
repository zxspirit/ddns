package provider

import (
	"fmt"
	"net/netip"
	"time"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"github.com/libdns/libdns"
)

func cfZone2LibdnsZone(zone zones.Zone) libdns.Zone {
	name := zone.Name
	name = libdns.AbsoluteName(name, "")
	return libdns.Zone{
		Name: name,
	}
}

func recordResponse2Record(response dns.RecordResponse) (libdns.Record, error) {

	switch response.Type {
	case "A", "AAAA":
		addr, err := netip.ParseAddr(response.Content)
		if err != nil {
			return nil, err
		}
		a := libdns.Address{
			Name: response.Name,
			TTL:  time.Duration(response.TTL) * time.Second,
			IP:   addr,
		}

		return a, nil
	case "CNAME":
		a := libdns.CNAME{
			Name:   response.Name,
			TTL:    time.Duration(response.TTL) * time.Second,
			Target: response.Content,
		}

		return a, nil
	case "TXT":
		txt := libdns.TXT{
			Name: response.Name,
			TTL:  time.Duration(response.TTL) * time.Second,
			Text: response.Content,
		}

		return txt, nil
	default:
		return nil, fmt.Errorf("unknown response type %s", response.Type)
	}
}
func zoneCache2RecordNewParams(recs libdns.Record) (dns.RecordNewParamsBodyUnion, error) {
	switch recs.RR().Type {
	case "A":
		return dns.ARecordParam{
			Name:    cloudflare.F(recs.RR().Name),
			TTL:     cloudflare.F(dns.TTL(recs.RR().TTL)),
			Type:    cloudflare.F(dns.ARecordTypeA),
			Content: cloudflare.F(recs.RR().Data),
			Proxied: cloudflare.F(false),
		}, nil
	case "CNAME":
		return dns.CNAMERecordParam{
			Name:    cloudflare.F(recs.RR().Name),
			TTL:     cloudflare.F(dns.TTL(recs.RR().TTL)),
			Type:    cloudflare.F(dns.CNAMERecordTypeCNAME),
			Content: cloudflare.F(recs.RR().Data),
			Proxied: cloudflare.F(false),
		}, nil
	case "TXT":
		return dns.TXTRecordParam{
			Name:    cloudflare.F(recs.RR().Name),
			TTL:     cloudflare.F(dns.TTL(recs.RR().TTL)),
			Type:    cloudflare.F(dns.TXTRecordTypeTXT),
			Content: cloudflare.F(recs.RR().Data),
			Proxied: cloudflare.F(false),
		}, nil
	case "AAAA":
		return dns.AAAARecordParam{
			Name:    cloudflare.F(recs.RR().Name),
			TTL:     cloudflare.F(dns.TTL(recs.RR().TTL)),
			Type:    cloudflare.F(dns.AAAARecordTypeAAAA),
			Content: cloudflare.F(recs.RR().Data),
			Proxied: cloudflare.F(false),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported record type %s", recs.RR().Type)
	}
}
func zoneCache2RecordUpdateParams(recs libdns.Record) (dns.RecordUpdateParamsBodyUnion, error) {
	switch recs.RR().Type {
	case "A":
		return dns.ARecordParam{
			Name:    cloudflare.F(recs.RR().Name),
			TTL:     cloudflare.F(dns.TTL(recs.RR().TTL)),
			Type:    cloudflare.F(dns.ARecordTypeA),
			Content: cloudflare.F(recs.RR().Data),
			Proxied: cloudflare.F(false),
		}, nil
	case "CNAME":
		return dns.CNAMERecordParam{
			Name:    cloudflare.F(recs.RR().Name),
			TTL:     cloudflare.F(dns.TTL(recs.RR().TTL)),
			Type:    cloudflare.F(dns.CNAMERecordTypeCNAME),
			Content: cloudflare.F(recs.RR().Data),
			Proxied: cloudflare.F(false),
		}, nil
	case "TXT":
		return dns.TXTRecordParam{
			Name:    cloudflare.F(recs.RR().Name),
			TTL:     cloudflare.F(dns.TTL(recs.RR().TTL)),
			Type:    cloudflare.F(dns.TXTRecordTypeTXT),
			Content: cloudflare.F(recs.RR().Data),
			Proxied: cloudflare.F(false),
		}, nil
	case "AAAA":
		return dns.AAAARecordParam{
			Name:    cloudflare.F(recs.RR().Name),
			TTL:     cloudflare.F(dns.TTL(recs.RR().TTL)),
			Type:    cloudflare.F(dns.AAAARecordTypeAAAA),
			Content: cloudflare.F(recs.RR().Data),
			Proxied: cloudflare.F(false),
		}, nil
	default:
		return nil, fmt.Errorf("unsupported record type %s", recs.RR().Type)
	}

}
