package pkg

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"log"
)

func ShowNextSteps(ctx context.Context, cf *cloudflare.API, settingsWithErrors []string, sourceZone *cloudflare.Zone) error {
	log.Printf("Next steps:")
	log.Printf(" - Check and update settings manualy:")
	for _, er := range settingsWithErrors {
		log.Printf("   - %s\n", er)
	}
	log.Printf(" - Migrate load balancers")
	lbs, err := cf.ListLoadBalancers(ctx, sourceZone.ID)
	if err != nil {
		return fmt.Errorf("can't load source load balancers: %v", err)
	}
	for _, lb := range lbs {
		log.Printf("   -> %s", lb.Name)
	}
	return nil
}