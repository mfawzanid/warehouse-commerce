package repository

import (
	"fmt"
	"mfawzanid/warehouse-commerce/core/entity"
)

func generateProductReservedKey(req *entity.LockOrderProductRequest) string {
	return fmt.Sprintf("reserved:%s:%s:%s", req.ProductId, req.WarehouseId, req.UserId)
}

func generateAllProductReservedKey(productId, warehouseId string) string {
	return fmt.Sprintf("reserved:%s:%s:*", productId, warehouseId)
}
