package pkg

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
)

func GetAccountZones(ctx context.Context, cf *cloudflare.API, account *cloudflare.Account, zoneName string) (*cloudflare.Zone, error) {
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
