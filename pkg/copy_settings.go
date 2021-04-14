package pkg

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"log"
)

func CopySettings(ctx context.Context, cf *cloudflare.API, sourceZone *cloudflare.Zone, destZone *cloudflare.Zone) ([]string, error) {
	log.Printf("Copy settings...")
	var settingsWithErrors []string
	allSettings, err := cf.ZoneSettings(ctx, sourceZone.ID)
	if err != nil {
		return nil, fmt.Errorf("can't get source zone settings: %v", err)
	}
	for _, s := range allSettings.Result {
		sets := []cloudflare.ZoneSetting{s}
		_, err = cf.UpdateZoneSettings(ctx, destZone.ID, sets)
		if err != nil {
			settingsWithErrors = append(settingsWithErrors, fmt.Sprintf("%s: %v", s.ID, err))
		} else {
			log.Printf(" - %s: %v", s.ID, s.Value)
		}
	}
	return settingsWithErrors, nil
}
