package main

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"log"
)

func CopyDNSRecords(sourceZone *cloudflare.Zone, destZone *cloudflare.Zone) error {
	log.Printf("Copy DNS records...")
	records, err := cf.DNSRecords(ctx, sourceZone.ID, cloudflare.DNSRecord{})
	if err != nil {
		return fmt.Errorf("can't get source zone records: %v", err)
	}
	for _, r := range records {
		log.Printf(" - create record %s %s...", r.Type, r.Name)
		_, err = cf.CreateDNSRecord(ctx, destZone.ID, r)
		if err != nil {
			return fmt.Errorf("can't create record: %v", err)
		}
	}
	return nil
}
