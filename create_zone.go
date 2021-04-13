package main

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"log"
)

func CreateZone(sourceZone *cloudflare.Zone, destAccount *cloudflare.Account, destZone *cloudflare.Zone) (*cloudflare.Zone, error) {
	log.Printf("\nCreate destination zone...")
	z, err := cf.CreateZone(ctx, sourceZone.Name, false, *destAccount, "full")
	if err != nil {
		return nil, fmt.Errorf("can't create destination zone details: %v", err)
	}
	destZone = &z
	return destZone, nil
}
