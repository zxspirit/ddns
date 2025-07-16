package local

import (
	"reflect"
	"testing"
)

func Test_getLocalPublicIpv6(t *testing.T) {
	tests := []struct {
		name      string
		wantIpv4s []string
		wantIpv6s []string
	}{
		{"t", []string{}, []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotIpv4s, gotIpv6s := getPublicIp()
			if !reflect.DeepEqual(gotIpv4s, tt.wantIpv4s) {
				t.Errorf("getPublicIp() gotIpv4s = %v, want %v", gotIpv4s, tt.wantIpv4s)
			}
			if !reflect.DeepEqual(gotIpv6s, tt.wantIpv6s) {
				t.Errorf("getPublicIp() gotIpv6s = %v, want %v", gotIpv6s, tt.wantIpv6s)
			}
		})
	}
}
