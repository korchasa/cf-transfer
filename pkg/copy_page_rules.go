package pkg

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"log"
)

func CopyPageRules(ctx context.Context, cf *cloudflare.API, sourceZone *cloudflare.Zone, destZone *cloudflare.Zone) error {
	cr, err := cf.ListPageRules(ctx, sourceZone.ID)
	if err != nil {
		return fmt.Errorf("can't list page rules: %v", err)
	}
	if len(cr) > 0 {
		log.Printf("Copy page rules...")
		for _, r := range cr {
			log.Printf(" - create page rule %v...", r.Targets)
			r.ID = ""
			_, err = cf.CreatePageRule(ctx, destZone.ID, r)
			if err != nil {
				return fmt.Errorf("can't create page rule: %v", err)
			}
		}
	}

	return nil
}
