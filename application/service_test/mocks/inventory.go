package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/dane4k/MerchShop/presentations/dto/response"
)

type InventoryRepo struct {
	mock.Mock
}

func (irm *InventoryRepo) GetUserInventory(ctx context.Context, userID int) ([]*response.InventoryItem, error) {
	args := irm.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*response.InventoryItem), args.Error(1)
}
