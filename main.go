package main

import (
	"context"
	"github.com/cloudflare/cloudflare-go"
	"log"
	"os"
)

var cf *cloudflare.API
var ctx context.Context

func main() {
	log.SetFlags(0)

	var err error
	cf, err = cloudflare.New(os.Getenv("CLOUDFLARE_KEY"), os.Getenv("CLOUDFLARE_EMAIL"))
	if err != nil {
		log.Fatalf("Can't init cloudflare cf: %s", err)
	}
	ctx = context.Background()

	sourceAccount, err := SelectAccount("Select source account")
	if err != nil {
		log.Fatalf("Can't select source account: %s", err)
	}

	sourceZone, err := SelectZone(sourceAccount, "Select source zone")
	if err != nil {
		log.Fatalf("Can't select source zone: %s", err)
	}

	destAccount, err := SelectAccount("Select destination account")
	if err != nil {
		log.Fatalf("Can't select source account: %s", err)
	}

	destZone, err := GetAccountZones(destAccount, sourceZone.Name)
	if err != nil {
		log.Fatalf("Can't find destination zone by name `%s`: %s", sourceZone.Name, err)
	}
	if destZone != nil {
		if err := CleanupZone(destZone); err != nil {
			log.Fatalf("Can't cleanup destination zone: %v",err)
		}
	} else {
		destZone, err = CreateZone(sourceZone, destAccount, destZone)
		if err != nil {
			log.Fatalf("Can't create zone: %s",err)
		}
	}

	if err = CopyDNSRecords(sourceZone, destZone); err != nil {
		log.Fatalf("Can't copy DNS records: %s", err)
	}

	settingsWithErrors, err := CopySettings(err, sourceZone, destZone)
	if err != nil {
		log.Fatalf("Can't copy zone settings: %s", err)
	}

	if err = CopyPageRules(sourceZone, destZone); err != nil {
		log.Fatalf("Can't copy page rules: %s", err)
	}

	if err := ShowNextSteps(settingsWithErrors, sourceZone); err != nil {
		log.Fatalf("Can't show next steps: %s", err)
	}
}
