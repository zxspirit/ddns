package provider

import "github.com/libdns/libdns"

type recordCache struct {
	libdns.Record
	ID string
}
