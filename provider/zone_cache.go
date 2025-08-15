package provider

import (
	"fmt"
	"sync"

	"github.com/libdns/libdns"
)

type zoneCache struct {
	ID string
	libdns.Zone
	records []*recordCache
	mx      sync.RWMutex
}

func (z *zoneCache) addRecordCache(recordId string, record libdns.Record) {
	z.mx.Lock()
	defer z.mx.Unlock()
	for _, cache := range z.records {
		if cache.ID == recordId {
			// Record already exists, no need to add it again
			return
		}
	}
	z.records = append(z.records, &recordCache{
		Record: record,
		ID:     recordId,
	})
}
func (z *zoneCache) getRecord(name string, dnsType string) (*recordCache, error) {
	z.mx.RLock()
	defer z.mx.RUnlock()
	for _, rec := range z.records {
		if rec.RR().Name == name && rec.RR().Type == dnsType {
			return rec, nil
		}
	}
	return nil, fmt.Errorf("record not found")
}
func (z *zoneCache) getAllRecords() []libdns.Record {
	z.mx.RLock()
	defer z.mx.RUnlock()
	records := make([]libdns.Record, 0, len(z.records))
	for _, rec := range z.records {
		records = append(records, rec.Record)
	}

	return records
}

func (z *zoneCache) deleteRecord(id string) {
	z.mx.Lock()
	defer z.mx.Unlock()
	for i, rec := range z.records {
		if rec.ID == id {
			z.records = append(z.records[:i], z.records[i+1:]...)
			return
		}
	}

}
