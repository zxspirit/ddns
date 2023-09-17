package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/cloudflare/cloudflare-go"
)

func main() {
	var CF_API_KEY string = os.Getenv("CF_API_KEY")
	var CF_API_EMAIL string = os.Getenv("CF_API_EMAIL")
	var CF_ZONE_NAME string = os.Getenv("CF_ZONE_NAME")
	var CF_RECORD_NAME string = os.Getenv("CF_RECORD_NAME")
	//获取本机公网ip
	a, err := net.InterfaceAddrs()
	if err != nil {
		panic(err)
	}
	var ip string //本机的公网ip
	for _, addr := range a {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil && ipnet.IP.IsGlobalUnicast() {
			// fmt.Println(ipnet.IP.String())
			ip = ipnet.IP.String()
		}
	}
	api, err := cloudflare.New(CF_API_KEY, CF_API_EMAIL)
	if err != nil {
		panic(err)
	}
	//获取域名id
	zoneID, err := api.ZoneIDByName(CF_ZONE_NAME)
	if err != nil {
		panic(err)
	}
	ctx := context.Background()

	//获取域名的dns记录
	records, _, err := api.ListDNSRecords(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.ListDNSRecordsParams{
		Type: "A",
		Name: CF_RECORD_NAME,
	})
	if err != nil {
		panic(err)
	}
	if len(records) == 0 {
		dnsRecord, err := api.CreateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.CreateDNSRecordParams{
			Type:    "A",
			Name:    CF_RECORD_NAME,
			Content: ip,
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("创建ipv4记录成功,当前ip:%s\n", dnsRecord.Content)
	} else {
		//更新dns记录
		dnsRecord, err := api.UpdateDNSRecord(ctx, cloudflare.ZoneIdentifier(zoneID), cloudflare.UpdateDNSRecordParams{
			Type:    "A",
			Name:    CF_RECORD_NAME,
			Content: ip,
			ID:      records[0].ID,
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("更新ipv4记录成功,当前ip:%s\n", dnsRecord.Content)
	}

}
