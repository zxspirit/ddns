package ddns

import (
	"context"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/zxspirit/ddns/ip"
	"github.com/zxspirit/ddns/provider"
)

func init() {
	caddy.RegisterModule(&DdnsPlugin{})
}

type DdnsPlugin struct {
	ipUtil ip.GetIp
	*provider.CloudflareProvider

	HostName string
}

func (d *DdnsPlugin) Cleanup() error {
	return nil
}

func (d *DdnsPlugin) Init() {
	d.ipUtil = &ip.LocalInterface{}
	cloudflareProvider := provider.NewCloudflareProvider()
	_ = cloudflareProvider.InitCache(context.Background())
	d.CloudflareProvider = cloudflareProvider
}

func (d *DdnsPlugin) Provision(c caddy.Context) error {

	return nil
}

func (d *DdnsPlugin) CaddyModule() caddy.ModuleInfo {

	return caddy.ModuleInfo{
		ID: "dns.providers.cf",
		New: func() caddy.Module {
			d2 := &DdnsPlugin{}
			d2.Init()
			return d2
		},
	}
}
func (d *DdnsPlugin) UnmarshalCaddyfile(dispenser *caddyfile.Dispenser) error {
	//if !dispenser.Next() {
	//	return fmt.Errorf("ddnsplugin: expected at least one argument")
	//}
	//if !dispenser.NextBlock(0) {
	//	return fmt.Errorf("ddnsplugin: expected at least one block")
	//}
	//if dispenser.Val() != "host_name" {
	//	return fmt.Errorf("ddnsplugin: expected 'host_name' directive, got '%s'", dispenser.Val())
	//}
	//if !dispenser.NextArg() {
	//	return fmt.Errorf("ddnsplugin: expected host name argument")
	//}
	//d.HostName = dispenser.Val()
	//if d.HostName == "" {
	//	return fmt.Errorf("ddnsplugin: host name cannot be empty")
	//}

	return nil
}

var (
	_ caddyfile.Unmarshaler = (*DdnsPlugin)(nil)
	_ caddy.Provisioner     = (*DdnsPlugin)(nil)
	_ caddy.CleanerUpper    = (*DdnsPlugin)(nil)
)
