package main

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/manifoldco/promptui"
	"log"
	"os"
)

func CleanupZone(destZone *cloudflare.Zone) error {
	log.Printf("Cleanup destination zone...")

	prompt := promptui.Select{
		Label: "Zone exists. Override?",
		Items: []string{"No", "Yes"},
	}
	_, result, err := prompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed %v", err)
	}
	if "Yes" != result {
		os.Exit(0)
	}

	log.Printf(" - Delete page rules...")
	prs, err := cf.ListPageRules(ctx, destZone.ID)
	if err != nil {
		return fmt.Errorf("can't get page rules: %v", err)
	}
	for _, pr := range prs {
		log.Printf("   - %s", pr.Targets)
		err := cf.DeletePageRule(ctx, destZone.ID, pr.ID)
		if err != nil {
			return fmt.Errorf("can't delete page rule: %v", err)
		}
	}

	log.Printf(" - Delete DNS records...")
	records, err := cf.DNSRecords(ctx, destZone.ID, cloudflare.DNSRecord{})
	if err != nil {
		return fmt.Errorf("can't get zone records: %v", err)
	}
	for _, r := range records {
		log.Printf("   - %s", r.Name)
		err = cf.DeleteDNSRecord(ctx, destZone.ID, r.ID)
		if err != nil {
			return fmt.Errorf("can't delete zone record: %v", err)
		}
	}
	return nil
}
