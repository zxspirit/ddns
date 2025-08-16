package provider

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/cloudflare/cloudflare-go/v4"
	"github.com/cloudflare/cloudflare-go/v4/dns"
	"github.com/cloudflare/cloudflare-go/v4/option"
	"github.com/cloudflare/cloudflare-go/v4/zones"
	"github.com/libdns/libdns"
)

var client = cloudflare.NewClient(option.WithAPIToken(os.Getenv("CLOUDFLARE_API_TOKEN")))

type CloudflareProvider struct {
	zoneCaches []*zoneCache
	mx         sync.RWMutex // Protects zoneCaches
}

func NewCloudflareProvider() *CloudflareProvider {
	return &CloudflareProvider{
		zoneCaches: make([]*zoneCache, 0),
	}
}
func (c *CloudflareProvider) InitCache(ctx context.Context) error {
	res, err := client.Zones.List(ctx, zones.ZoneListParams{Status: cloudflare.F(zones.ZoneListParamsStatusActive)})
	if err != nil {
		return err
	}
	for _, zone := range res.Result {
		libdnsZone := cfZone2LibdnsZone(zone)
		zoneC, err := c.addZoneCache(zone.ID, libdnsZone)
		if err != nil {
			return err
		}
		list, err := client.DNS.Records.List(ctx, dns.RecordListParams{
			ZoneID: cloudflare.F(zone.ID),
		})
		if err != nil {
			return err
		}
		for _, response := range list.Result {
			record, err := recordResponse2Record(response)
			if err != nil {
				return err
			}
			zoneC.addRecordCache(response.ID, record)
		}
	}
	return nil
}

func (c *CloudflareProvider) addZoneCache(id string, zone libdns.Zone) (*zoneCache, error) {
	c.mx.Lock()
	defer c.mx.Unlock()
	// Check if the zone already exists in the cache
	for _, existingZone := range c.zoneCaches {
		if existingZone.ID == id {
			// Zone already exists, no need to add it again
			return nil, fmt.Errorf("zone %s already exists", id)
		}
	}
	r := &zoneCache{
		ID:      id,
		Zone:    zone,
		records: make([]*recordCache, 0),
	}
	c.zoneCaches = append(c.zoneCaches, r)
	return r, nil
}

func (c *CloudflareProvider) getZoneCacheByName(name string) (*zoneCache, error) {
	c.mx.RLock()
	defer c.mx.RUnlock()
	for _, zone := range c.zoneCaches {
		absoluteName := libdns.AbsoluteName(name, "")
		if zone.Name == absoluteName {
			return zone, nil
		}
	}
	return nil, fmt.Errorf("zone %s not found", name)
}
func (c *CloudflareProvider) getAllZones() []libdns.Zone {
	c.mx.RLock()
	defer c.mx.RUnlock()
	r := make([]libdns.Zone, 0, len(c.zoneCaches))
	for _, cache := range c.zoneCaches {
		r = append(r, cache.Zone)
	}
	return r
}

func (c *CloudflareProvider) ListZones(_ context.Context) ([]libdns.Zone, error) {
	return c.getAllZones(), nil
}

func (c *CloudflareProvider) GetRecords(_ context.Context, zone string) ([]libdns.Record, error) {
	zC, err := c.getZoneCacheByName(zone)
	if err != nil {
		return nil, err
	}
	return zC.getAllRecords(), nil
}

// SetRecords 如果没有记录则创建，如果有记录但内容不同则更新
// 如果有记录且内容相同则不做任何操作
// 返回所有操作的记录
// 如果有记录但内容相同则不返回该记录
// 返回的记录是更新后的记录
func (c *CloudflareProvider) SetRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	z, err := c.getZoneCacheByName(zone)
	if err != nil {
		return nil, err
	}
	res := make([]libdns.Record, 0, len(recs))
	for _, rec := range recs {
		record, err := z.getRecord(rec.RR().Name, rec.RR().Type)
		if err != nil {
			params, err2 := zoneCache2RecordNewParams(rec)
			if err2 != nil {
				return nil, err2
			}
			r, err := client.DNS.Records.New(ctx, dns.RecordNewParams{
				ZoneID: cloudflare.F(z.ID),
				Body:   params,
			})
			if err != nil {
				continue
			}
			response2Record, err := recordResponse2Record(*r)
			if err != nil {
				continue
			}
			res = append(res, response2Record)
			z.addRecordCache(r.ID, response2Record)

		} else {
			if record.RR().Data != rec.RR().Data {
				params, err := zoneCache2RecordUpdateParams(rec)
				if err != nil {
					return nil, fmt.Errorf("failed to convert record to params: %w", err)
				}
				update, err := client.DNS.Records.Update(ctx, record.ID, dns.RecordUpdateParams{
					ZoneID: cloudflare.F(z.ID),
					Body:   params,
				})
				if err != nil {
					continue
				}
				response2Record, err := recordResponse2Record(*update)
				if err != nil {
					continue
				}
				res = append(res, response2Record)
				z.addRecordCache(update.ID, response2Record)
			}
		}
	}
	return res, nil
}

// AppendRecords 如果没有记录则创建，如果有记录不做任何操作
// 返回所有操作成功的记录
func (c *CloudflareProvider) AppendRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {
	z, err := c.getZoneCacheByName(zone)
	if err != nil {
		return nil, err
	}
	res := make([]libdns.Record, 0, len(recs))
	for _, input := range recs {
		name := input.RR().Name
		name = libdns.AbsoluteName(name, z.Name)
		_, err := z.getRecord(name, input.RR().Type)
		if err == nil {
			// Record already exists, skip it
			continue
		}
		params, err := zoneCache2RecordNewParams(input)
		if err != nil {
			return nil, err
		}
		r, err := client.DNS.Records.New(ctx, dns.RecordNewParams{
			ZoneID: cloudflare.F(z.ID),
			Body:   params,
		},
		)
		if err != nil {
			continue
		}
		response2Record, err := recordResponse2Record(*r)
		if err != nil {
			return nil, fmt.Errorf("failed to convert response to record: %w", err)
		}
		z.addRecordCache(r.ID, response2Record)
		// Append the new record to the result
		res = append(res, response2Record)
	}

	return res, nil
}

func (c *CloudflareProvider) DeleteRecords(ctx context.Context, zone string, recs []libdns.Record) ([]libdns.Record, error) {

	zc, err := c.getZoneCacheByName(zone)
	if err != nil {
		return nil, err
	}
	res := make([]libdns.Record, 0, len(recs))
	for _, rec := range recs {
		record, err := zc.getRecord(rec.RR().Name, rec.RR().Type)
		if err != nil {
			continue
		}
		_, err = client.DNS.Records.Delete(ctx, record.ID, dns.RecordDeleteParams{ZoneID: cloudflare.F(zc.ID)})
		if err != nil {
			return nil, err
		}
		res = append(res, record)
		zc.deleteRecord(record.ID)
	}

	return res, nil

}

var (
	_ libdns.RecordGetter   = (*CloudflareProvider)(nil)
	_ libdns.RecordAppender = (*CloudflareProvider)(nil)
	_ libdns.RecordSetter   = (*CloudflareProvider)(nil)
	_ libdns.RecordDeleter  = (*CloudflareProvider)(nil)
	_ libdns.ZoneLister     = (*CloudflareProvider)(nil)
)
