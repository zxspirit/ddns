package ddnsplugin

import (
	"context"
	"fmt"
	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/libdns/libdns"
	"github.com/sirupsen/logrus"
	"github.com/zxspirit/cflibdns"
	"github.com/zxspirit/ddns/ip"
	"strings"
	"time"
)

func init() {
	caddy.RegisterModule(&DdnsPlugin{})
}

type DdnsPlugin struct {
	ipUtil ip.GetIp
	*cflibdns.Provider
	HostName string
}

func (d *DdnsPlugin) Provision(c caddy.Context) error {
	initCache := func() {
		err := d.InitCache(c)
		if err != nil {
			logrus.Errorf("error initializing cache: %v", err)
			return
		}
		ticker := time.NewTicker(24 * time.Hour)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				err := d.InitCache(c)
				if err != nil {
					logrus.Errorf("error initializing cache: %v", err)
					continue
				}
			case <-c.Done():
				return

			}
		}
	}
	setHostName := func(domain string) {
		if domain == "" {
			logrus.Error("ddnsplugin: host name is empty")
			return
		}
		if err := d.ddns(c, domain); err != nil {
			logrus.Errorf("ddns update failed: %v", err)
			return
		}
		newTicker := time.NewTicker(30 * time.Second)
		defer newTicker.Stop()
		for {
			select {
			case <-newTicker.C:
				if err := d.ddns(c, domain); err != nil {
					logrus.Errorf("ddns update failed: %v", err)
					continue
				}
			case <-c.Done():
				return
			}
		}
	}
	var domains []string
	domains = append(domains, "aaaa.newzhxu.com", "abcc.newzhxu.com")
	setServers := func(domains []string) {
		err := d.server2Acme(c, domains)
		if err != nil {
			logrus.Errorf("error setting servers for ACME: %v", err)
			return
		}
		newTicker := time.NewTicker(30 * time.Second)
		defer newTicker.Stop()
		for {
			select {
			case <-newTicker.C:
				if err := d.server2Acme(c, domains); err != nil {
					logrus.Errorf("error setting servers for ACME: %v", err)
					continue
				}
			case <-c.Done():
				return
			}
		}

	}
	go initCache()
	go setHostName(d.HostName)
	go setServers(domains)

	return nil
}

func (d *DdnsPlugin) CaddyModule() caddy.ModuleInfo {

	return caddy.ModuleInfo{
		ID: "dns.providers.ddnsplugin",
		New: func() caddy.Module {
			p := cflibdns.New(logrus.New())
			return &DdnsPlugin{
				ipUtil:   ip.LocalInterface{},
				Provider: p,
			}
		},
	}
}
func (d *DdnsPlugin) UnmarshalCaddyfile(dispenser *caddyfile.Dispenser) error {
	if !dispenser.Next() {
		return fmt.Errorf("ddnsplugin: expected at least one argument")
	}
	if !dispenser.NextBlock(0) {
		return fmt.Errorf("ddnsplugin: expected at least one block")
	}
	if dispenser.Val() != "host_name" {
		return fmt.Errorf("ddnsplugin: expected 'host_name' directive, got '%s'", dispenser.Val())
	}
	if !dispenser.NextArg() {
		return fmt.Errorf("ddnsplugin: expected host name argument")
	}
	d.HostName = dispenser.Val()
	if d.HostName == "" {
		return fmt.Errorf("ddnsplugin: host name cannot be empty")
	}

	return nil
}
func belongsToZone(domain string, zone string) bool {
	z := strings.ToLower(zone)
	d := strings.ToLower(domain)
	if strings.HasSuffix(d, ".") {
		d = d[:len(d)-1]
	}
	if strings.HasSuffix(z, ".") {
		z = z[:len(z)-1]
	}
	return d == z || strings.HasSuffix(d, "."+z)
}

func (d *DdnsPlugin) server2Acme(ctx context.Context, domains []string) error {
	zones, err := d.ListZones(ctx)
	if err != nil {
		return fmt.Errorf("error listing zones: %v", err)
	}
	for _, zone := range zones {
		for _, domain := range domains {
			record := fmt.Sprintf("%s.%s", d.HostName, zone.Name)
			if belongsToZone(domain, zone.Name) {
				var toDelete []libdns.Record
				a := libdns.RR{
					Name: domain,
					TTL:  1,
					Type: "A",
					Data: "",
				}
				aaaa := libdns.RR{
					Name: domain,
					TTL:  1,
					Type: "AAAA",
					Data: "",
				}
				toDelete = append(toDelete, a, aaaa)
				_, err := d.Provider.DeleteRecords(ctx, zone.Name, toDelete)
				if err != nil {
					return fmt.Errorf("error deleting records for %s: %w", domain, err)
				}

				var recs []libdns.Record
				recs = append(recs, libdns.CNAME{
					Name:   domain,
					TTL:    1,
					Target: record,
				})
				_, err = d.Provider.SetRecords(ctx, zone.Name, recs)
				if err != nil {
					return fmt.Errorf("error setting records for %s: %w", domain, err)
				}

			}
		}
	}
	return nil

}

func (d *DdnsPlugin) ddns(ctx context.Context, hostname string) error {
	zones, err := d.ListZones(ctx)
	if err != nil {
		return fmt.Errorf("error listing zones: %w", err)
	}
	ipv4, ipv6, err := d.ipUtil.GetPublicIp()
	if err != nil {
		return fmt.Errorf("error getting public IP: %w", err)
	}

	for _, zone := range zones {
		domain := fmt.Sprintf("%s.%s", hostname, zone.Name)
		records := make([]libdns.Record, 0, 2)
		v4 := libdns.RR{
			Name: domain,
			TTL:  1,
			Type: "A",
			Data: ipv4,
		}
		v6 := libdns.RR{
			Name: domain,
			TTL:  1,
			Type: "AAAA",
			Data: ipv6,
		}
		records = append(records, v4, v6)
		_, err = d.SetRecords(ctx, zone.Name, records)
		if err != nil {
			return fmt.Errorf("error setting records for %s: %w", domain, err)
		}

	}
	return nil
}

var (
	_ caddyfile.Unmarshaler = (*DdnsPlugin)(nil)
	_ caddy.Provisioner     = (*DdnsPlugin)(nil)
)
