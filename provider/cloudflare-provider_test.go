package provider

import (
	"context"
	"reflect"
	"testing"

	"github.com/libdns/libdns"
)

func TestCloudflareProvider_ListZones(t *testing.T) {
	type fields struct {
		zoneCaches []*zoneCache
	}
	type args struct {
		in0 context.Context
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []libdns.Zone
		wantErr bool
	}{
		{"t1", fields{
			zoneCaches: make([]*zoneCache, 0),
		}, args{in0: context.Background()}, make([]libdns.Zone, 0), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := NewCloudflareProvider()
			err2 := c.InitCache(tt.args.in0)
			if err2 != nil {
				t.Errorf("InitCache() error = %v", err2)
				return
			}
			got, err := c.ListZones(tt.args.in0)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListZones() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ListZones() got = %v, want %v", got, tt.want)
			}
		})
	}
}
