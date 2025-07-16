package ddnsplugin

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/zxspirit/cflibdns"
	"github.com/zxspirit/ddns/ip"
	"testing"
)

func TestDdnsPlugin_ddns(t *testing.T) {
	type fields struct {
		ipUtil   ip.GetIp
		Provider *cflibdns.Provider
	}
	type args struct {
		hostname string
	}
	provider := cflibdns.New(logrus.New())
	err := provider.InitCache(context.Background())
	if err != nil {
		t.Fatalf("Failed to initialize cache: %v", err)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"t1", fields{ipUtil: ip.LocalInterface{}, Provider: provider}, args{hostname: "test"}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := DdnsPlugin{
				ipUtil:   tt.fields.ipUtil,
				Provider: tt.fields.Provider,
			}
			if err := d.ddns(context.Background(), tt.args.hostname); (err != nil) != tt.wantErr {
				t.Errorf("ddns() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_belongsToZone(t *testing.T) {
	type args struct {
		domain string
		zone   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"t1", args{domain: "example.com", zone: "example.com"}, true},
		{"t2", args{domain: "sub.example.com.", zone: "example.com."}, true},
		{"t3", args{domain: "example.cOm.", zone: "example.Com"}, true},
		{"t4", args{domain: "example.com.", zone: "example.co"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := belongsToZone(tt.args.domain, tt.args.zone); got != tt.want {
				t.Errorf("belongsToZone() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDdnsPlugin_server2Acme(t *testing.T) {
	type fields struct {
		ipUtil   ip.GetIp
		Provider *cflibdns.Provider
		HostName string
	}
	type args struct {
		ctx     context.Context
		domains []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{"t1", fields{
			ipUtil:   ip.LocalInterface{},
			Provider: cflibdns.New(logrus.New()),
			HostName: "test",
		}, args{
			ctx:     context.Background(),
			domains: []string{"a.newzhxu.com"},
		}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			provider := tt.fields.Provider
			err := provider.InitCache(tt.args.ctx)
			if err != nil {
				t.Fatalf("Failed to initialize cache: %v", err)
			}
			d := &DdnsPlugin{
				ipUtil:   tt.fields.ipUtil,
				Provider: provider,
				HostName: tt.fields.HostName,
			}
			if err := d.server2Acme(tt.args.ctx, tt.args.domains); (err != nil) != tt.wantErr {
				t.Errorf("server2Acme() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
