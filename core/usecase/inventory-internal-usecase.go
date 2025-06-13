package usecase

import (
	"fmt"
	"mfawzanid/warehouse-commerce/core/entity"
	errorutil "mfawzanid/warehouse-commerce/utils/error"
)

func (u *inventoryUsecase) getWarehouseId(warehouseName string) (string, error) {
	getWarehousesResp, err := u.inventoryRepo.GetWarehouses(&entity.GetWarehousesRequest{
		Name: warehouseName,
	})
	if err != nil {
		return "", err
	}

	if len(getWarehousesResp.Warehouses) > 0 {
		return getWarehousesResp.Warehouses[0].Id, nil
	}

	return "", errorutil.NewErrorCode(errorutil.ErrNotFound, fmt.Errorf("error get warehouse by name '%v'", warehouseName))
}

func (u *inventoryUsecase) getShopId(shopName string) (string, error) {
	getShopsResp, err := u.inventoryRepo.GetShops(&entity.GetShopsRequest{
		Name: shopName,
	})
	if err != nil {
		return "", err
	}

	if len(getShopsResp.Shops) > 0 {
		return getShopsResp.Shops[0].Id, nil
	}

	return "", errorutil.NewErrorCode(errorutil.ErrNotFound, fmt.Errorf("error get shop by name '%v'", shopName))
}

func (u *inventoryUsecase) getProductId(productName string) (string, error) {
	product, err := u.inventoryRepo.GetProductByName(productName)
	if err != nil {
		return "", err
	}

	return product.Id, nil
}
