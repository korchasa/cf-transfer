package main

import (
	"fmt"
	"github.com/cloudflare/cloudflare-go"
)

func GetAccountZones(account *cloudflare.Account, zoneName string) (*cloudflare.Zone, error) {
	zl, err := cf.ListZones(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get zones list: %v", err)
	}
	for _, z := range zl {
		if z.Account.ID == account.ID && z.Name == zoneName {
			return &z, nil
		}
	}
	return nil, nil
}
