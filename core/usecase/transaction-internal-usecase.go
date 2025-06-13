package usecase

import (
	"context"
	"fmt"
	"mfawzanid/warehouse-commerce/core/entity"
	errorutil "mfawzanid/warehouse-commerce/utils/error"
)

func (u *transactionUsecase) validatePaymentAmount(orderId string, paymentAmount int) error {
	var isActiveOrder *bool
	active := true
	isActiveOrder = &active

	if orderId == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error pay order: order id is mandatory"))
	}

	order, err := u.transactionRepo.GetOrderById(orderId, isActiveOrder)
	if err != nil {
		return err
	}

	if paymentAmount != order.Amount {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error pay order: amount is not same with order's amount"))
	}

	return nil
}

func (u *transactionUsecase) getProductDetailMap(req *entity.OrderProductsRequest) (map[string]*entity.ProductDetail, error) {
	var productIds []string
	for _, item := range req.Items {
		productIds = append(productIds, item.ProductId)
	}

	getProductDetailsResp, err := u.inventoryRepo.GetProductDetailsByShopId(&entity.GetProductDetailsByShopIdRequest{
		ShopId:     req.ShopId,
		ProductIds: productIds,
	})
	if err != nil {
		return nil, err
	}
	if len(getProductDetailsResp.ProductDetails) < len(req.Items) {
		return nil, errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error order products: some products are not found"))
	}

	productDetailMap := make(map[string]*entity.ProductDetail)
	for _, productDetail := range getProductDetailsResp.ProductDetails {
		productDetailMap[productDetail.ProductId] = productDetail
	}

	return productDetailMap, nil
}

func (u *transactionUsecase) validateTotalStock(productDetail *entity.ProductDetail, item *entity.OrderProductItem) error {
	reservedQuantity, err := u.redisRepo.GetReservedProductQuantity(context.Background(), productDetail.ProductId, productDetail.WarehouseId)
	if err != nil {
		return err
	}

	remainingStock := productDetail.TotalStock - reservedQuantity - item.Quantity
	if remainingStock < 0 {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error order products: product '%s' stock is not sufficient", item.ProductId))
	}

	return nil
}
