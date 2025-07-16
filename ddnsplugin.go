package ddnsplugin

import (
	"github.com/caddyserver/caddy/v2"
	provider "github.com/zxspirit/cflibdns"
)

func init() {
	caddy.RegisterModule(DdnsPlugin{})
}

type DdnsPlugin struct {
	*provider.Provider
}

func (d DdnsPlugin) CaddyModule() caddy.ModuleInfo {

	return caddy.ModuleInfo{
		ID:  "dns.providers.ddnsplugin",
		New: func() caddy.Module { return &DdnsPlugin{provider.New()} },
	}
}
