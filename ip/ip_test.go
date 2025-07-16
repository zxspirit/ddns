package ip

import (
	"testing"
)

func TestLocalInterface_GetPublicIp(t *testing.T) {
	tests := []struct {
		name     string
		wantIpv4 string
		wantIpv6 string
		wantErr  bool
	}{
		{"t1", "", "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := LocalInterface{}
			gotIpv4, gotIpv6, err := l.GetPublicIp()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPublicIp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if gotIpv4 != tt.wantIpv4 {
				t.Errorf("GetPublicIp() gotIpv4 = %v, want %v", gotIpv4, tt.wantIpv4)
			}
			if gotIpv6 != tt.wantIpv6 {
				t.Errorf("GetPublicIp() gotIpv6 = %v, want %v", gotIpv6, tt.wantIpv6)
			}
		})
	}
}
