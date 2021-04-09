package main

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"log"
	"os"
)

var (
	settingsToSkip = []string{"ciphers", "cname_flattening", "http2",
		"log_to_cloudflare", "mirage", "mobile_redirect", "orange_to_orange",
		"origin_error_page_pass_thru", "polish", "prefetch_preload",
		"response_buffering", "sort_query_string_for_cache",
		"true_client_ip_header", "visitor_ip", "waf", "webp"}
)

func main() {

	log.SetFlags(0)

	if len(os.Args) != 3 {
		fmt.Printf("Usage: %s <source_zone_id> <dest_zone_id>\n", os.Args[0])
		os.Exit(1)
	}

	fromId := os.Args[1]
	toId := os.Args[2]

	api, err := cloudflare.New(os.Getenv("CLOUDFLARE_KEY"), os.Getenv("CLOUDFLARE_EMAIL"))
	if err != nil {
		log.Fatalf("Can't init cloudflare api: %s", err)
	}

	fromZone, err := api.ZoneDetails(context.Background(), fromId)
	if err != nil {
		log.Fatalf("Can't get source zone details: %s", err)
	}
	toZone, err := api.ZoneDetails(context.Background(), toId)
	if err != nil {
		log.Fatalf("Can't get destination zone details: %s", err)
	}
	log.Printf("Migrate records from `%s` to `%s`", fromZone.Name, toZone.Name)


	log.Printf("Copy settings...")
	allSettings, err := api.ZoneSettings(context.Background(), fromId)
	if err != nil {
		log.Fatalf("Can't get source zone settings: %s", err)
	}
ROOT:
	for _, s := range allSettings.Result {
		for _, ss := range settingsToSkip {
			if ss == s.ID {
				continue ROOT
			}
		}
		sets := []cloudflare.ZoneSetting{s}
		_, err = api.UpdateZoneSettings(context.Background(), toId, sets)
		if err != nil {
			log.Printf(" - %s: %v\t!!! %v", s.ID, s.Value, err)
		} else {
			log.Printf(" - %s: %v", s.ID, s.Value)
		}

	}

	log.Printf("Cleanup old DNS records...")
	records, err := api.DNSRecords(context.Background(), toId, cloudflare.DNSRecord{})
	if err != nil {
		log.Fatalf("Can't get destination zone records: %s", err)
	}
	for _, r := range records {
		log.Printf(" - delete record %s %s...", r.Type, r.Name)
		err = api.DeleteDNSRecord(context.Background(), toId, r.ID)
		if err != nil {
			log.Fatalf("Can't delete record: %s", err)
		}
	}

	log.Printf("Copy DNS records...")
	records, err = api.DNSRecords(context.Background(), fromId, cloudflare.DNSRecord{})
	if err != nil {
		log.Fatalf("Can't get source zone records: %s", err)
	}
	for _, r := range records {
		log.Printf(" - create record %s %s...", r.Type, r.Name)
		_, err = api.CreateDNSRecord(context.Background(), toId, r)
		if err != nil {
			log.Fatalf("Can't create record: %s", err)
		}
	}

	log.Printf("Cleanup old page rules...")
	dr, err := api.ListPageRules(context.Background(), toId)
	for _, r := range dr {
		log.Printf(" - delete page rule %v...", r.Targets)
		err = api.DeletePageRule(context.Background(), toId, r.ID)
		if err != nil {
			log.Fatalf("Can't delete page rule: %s", err)
		}
	}

	log.Printf("Copy page rules...")
	cr, err := api.ListPageRules(context.Background(), fromId)
	for _, r := range cr {
		log.Printf(" - create page rule %v...", r.Targets)
		r.ID = ""
		_, err = api.CreatePageRule(context.Background(), toId, r)
		if err != nil {
			log.Fatalf("Can't create page rule: %s", err)
		}
	}
}

