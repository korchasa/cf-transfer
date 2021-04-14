package pkg

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"log"
)

func ShowNextSteps(ctx context.Context, cf *cloudflare.API, settingsWithErrors []string, sourceZone *cloudflare.Zone) error {

	lbs, err := cf.ListLoadBalancers(ctx, sourceZone.ID)
	if err != nil {
		return fmt.Errorf("can't load source load balancers: %v", err)
	}

	if len(lbs) == 0 && len(settingsWithErrors) == 0 {
		return nil
	}

	log.Printf("Next steps:")

	if len(lbs) > 0 {
		log.Printf(" - Migrate load balancers")
		for _, lb := range lbs {
			log.Printf("   -> %s", lb.Name)
		}
	}

	if len(settingsWithErrors) > 0 {
		log.Printf(" - Check and update settings manualy:")
		for _, er := range settingsWithErrors {
			log.Printf("   - %s\n", er)
		}
	}

	return nil
}
