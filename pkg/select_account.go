package pkg

import (
	"context"
	"fmt"
	"github.com/cloudflare/cloudflare-go"
	"github.com/manifoldco/promptui"
)

func SelectAccount(ctx context.Context, cf *cloudflare.API, label string) (*cloudflare.Account, error) {
	acs, _, err := cf.Accounts(ctx, cloudflare.PaginationOptions{PerPage: 50})
	if err != nil {
		return nil, fmt.Errorf("can't get accounts list: %v", err)
	}

	itemsMap := make(map[string]cloudflare.Account)
	items := make([]string, 0)
	for _, ac := range acs {
		itemsMap[ac.Name] = ac
		items = append(items, ac.Name)
	}

	prompt := promptui.Select{
		Label: label,
		Items: items,
		Size: 20,
	}

	_, name, pErr := prompt.Run()
	if pErr != nil {
		return nil, fmt.Errorf("account prompt fail: %v", pErr)
	}

	a := itemsMap[name]
	return &a, nil
}
