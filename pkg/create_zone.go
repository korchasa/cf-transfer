package pkg

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"log"
)

func CreateZone(ctx context.Context, cf *cloudflare.API, sourceZone *cloudflare.Zone, destAccount *cloudflare.Account) (*cloudflare.Zone, error) {
	log.Printf("Create destination zone...")
	z, err := cf.CreateZone(ctx, sourceZone.Name, false, *destAccount, "full")
	if err != nil {
		return nil, fmt.Errorf("can't create destination zone details: %v", err)
	}
	destZone := &z
	return destZone, nil
}
