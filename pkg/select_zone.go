package pkg

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/manifoldco/promptui"
)

func SelectZone(ctx context.Context, cf *cloudflare.API, account *cloudflare.Account, label string) (*cloudflare.Zone, error) {
	zl, err := cf.ListZones(ctx)
	if err != nil {
		return nil, fmt.Errorf("can't get zones list: %v", err)
	}

	itemsMap := make(map[string]cloudflare.Zone)
	items := make([]string, 0)
	for _, z := range zl {
		if z.Account.ID != account.ID {
			continue
		}
		itemsMap[z.Name] = z
		items = append(items, z.Name)
	}

	prompt := promptui.Select{
		Label: label,
		Items: items,
		Size: 20,
	}

	_, name, pErr := prompt.Run()
	if pErr != nil {
		return nil, fmt.Errorf("zone prompt fail: %v", pErr)
	}

	za := itemsMap[name]
	return &za, nil
}
