package mocks

import (
	"context"

	"github.com/dane4k/MerchShop/internal/presentation/dto/response"
	"github.com/stretchr/testify/mock"
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
