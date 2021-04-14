package main

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	"korchasa/cf-transfer/pkg"
	"log"
	"os"
)

func init() {
	log.SetFlags(0)
}

func main() {
	var (
		cf  *cloudflare.API
		ctx context.Context
		err error
	)

	cf, err = cloudflare.New(os.Getenv("CLOUDFLARE_KEY"), os.Getenv("CLOUDFLARE_EMAIL"))
	if err != nil {
		log.Fatalf("Can't init cloudflare cf: %s", err)
	}
	ctx = context.Background()

	sourceAccount, err := pkg.SelectAccount(ctx, cf, "Select source account")
	if err != nil {
		log.Fatalf("Can't select source account: %s", err)
	}

	sourceZone, err := pkg.SelectZone(ctx, cf, sourceAccount, "Select source zone")
	if err != nil {
		log.Fatalf("Can't select source zone: %s", err)
	}

	destAccount, err := pkg.SelectAccount(ctx, cf, "Select destination account")
	if err != nil {
		log.Fatalf("Can't select source account: %s", err)
	}

	destZone, err := pkg.GetAccountZones(ctx, cf, destAccount, sourceZone.Name)
	if err != nil {
		log.Fatalf("Can't find destination zone by name `%s`: %s", sourceZone.Name, err)
	}
	if destZone != nil {
		if err := pkg.CleanupZone(ctx, cf, destZone); err != nil {
			log.Fatalf("Can't cleanup destination zone: %v", err)
		}
	} else {
		destZone, err = pkg.CreateZone(ctx, cf, sourceZone, destAccount)
		if err != nil {
			log.Fatalf("Can't create zone: %s", err)
		}
	}

	if err = pkg.CopyDNSRecords(ctx, cf, sourceZone, destZone); err != nil {
		log.Fatalf("Can't copy DNS records: %s", err)
	}

	settingsWithErrors, err := pkg.CopySettings(ctx, cf, sourceZone, destZone)
	if err != nil {
		log.Fatalf("Can't copy zone settings: %s", err)
	}

	if err = pkg.CopyPageRules(ctx, cf, sourceZone, destZone); err != nil {
		log.Fatalf("Can't copy page rules: %s", err)
	}

	if err := pkg.ShowNextSteps(ctx, cf, settingsWithErrors, sourceZone); err != nil {
		log.Fatalf("Can't show next steps: %s", err)
	}
}
