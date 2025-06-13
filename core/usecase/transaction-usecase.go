package usecase

import (
	"context"
	"fmt"
	"log"
	"mfawzanid/warehouse-commerce/core/entity"
	"mfawzanid/warehouse-commerce/core/repository"
	transactionutil "mfawzanid/warehouse-commerce/utils/transaction"
	"time"

	"github.com/gofrs/uuid/v5"
)

type TransactionUsecaseInterface interface {
	GetProductDetailsByShopId(req *entity.GetProductDetailsByShopIdRequest) (*entity.GetProductDetailsByShopIdResponse, error)
	OrderProducts(req *entity.OrderProductsRequest) (string, error)
	PayOrder(req *entity.PayOrderRequest) error
}

type transactionUsecase struct {
	inventoryRepo   repository.InventoryRepositoryInterface
	transactionRepo repository.TransactionRepositoryInterface
	redisRepo       repository.RedisRepositoryInterface
}

func NewTransactionUsecase(inventoryRepo repository.InventoryRepositoryInterface, transactionRepo repository.TransactionRepositoryInterface, redisRepo repository.RedisRepositoryInterface) TransactionUsecaseInterface {
	return &transactionUsecase{
		inventoryRepo:   inventoryRepo,
		transactionRepo: transactionRepo,
		redisRepo:       redisRepo,
	}
}

func (u *transactionUsecase) GetProductDetailsByShopId(req *entity.GetProductDetailsByShopIdRequest) (*entity.GetProductDetailsByShopIdResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	resp, err := u.inventoryRepo.GetProductDetailsByShopId(req)
	if err != nil {
		return nil, err
	}

	// check reserved quantity, if any then should return the real remaining stock
	for _, productDetail := range resp.ProductDetails {
		reservedQuantity, err := u.redisRepo.GetReservedProductQuantity(context.Background(), productDetail.ProductId, productDetail.WarehouseId)
		if err != nil {
			return nil, err
		}
		productDetail.TotalStock -= reservedQuantity
	}

	return resp, nil
}

func (u *transactionUsecase) OrderProducts(req *entity.OrderProductsRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}

	productDetailMap, err := u.getProductDetailMap(req)
	if err != nil {
		return "", err
	}

	newUUID, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("error order products in generating uuid: %v", err.Error())
	}
	orderId := newUUID.String()

	var amount int
	var orderItems []*entity.OrderItem

	for _, item := range req.Items {
		// validate stock for each product
		productDetail, ok := productDetailMap[item.ProductId]
		if !ok {
			return "", fmt.Errorf("error order products: product '%s' is not found", item.ProductId)
		}

		if err := u.validateTotalStock(productDetail, item); err != nil {
			return "", err
		}

		// lock the product and its quantity
		if err = u.redisRepo.LockOrderProduct(context.Background(), &entity.LockOrderProductRequest{
			ProductId:   productDetail.ProductId,
			WarehouseId: productDetail.WarehouseId,
			UserId:      req.UserId,
			Quantity:    item.Quantity,
		}); err != nil {
			return "", fmt.Errorf("error order product '%s': %v", item.ProductId, err.Error())
		}

		amount += productDetail.Price * item.Quantity

		orderItems = append(orderItems, &entity.OrderItem{
			OrderId:     orderId,
			ProductId:   item.ProductId,
			ShopId:      req.ShopId,
			WarehouseId: productDetail.WarehouseId,
			Quantity:    item.Quantity,
			UnitPrice:   productDetail.Price,
		})
	}

	timeNow := time.Now()
	if err := u.transactionRepo.InsertOrder(&entity.Order{
		Id:        orderId,
		ShopId:    req.ShopId,
		UserId:    req.UserId,
		Amount:    amount,
		Status:    entity.OrderStatusPending,
		CreatedAt: timeNow,
		ExpiredAt: timeNow.Add(time.Duration(entity.OrderExpireTimeInMinute * time.Minute)),
	}); err != nil {
		return "", err
	}

	if err := u.transactionRepo.InsertOrderItems(orderItems); err != nil {
		return "", err
	}

	return orderId, nil
}

func (u *transactionUsecase) PayOrder(req *entity.PayOrderRequest) error {
	if err := u.validatePaymentAmount(req.OrderId, req.Amount); err != nil {
		return err
	}

	// insert payment & update order are execute in one transaction
	tx, err := u.inventoryRepo.GetDb().Begin()
	if err != nil {
		return fmt.Errorf("error pay order in initiating transaction: %v", err.Error())
	}

	defer func() {
		err = transactionutil.SettleTransaction(tx, err)
	}()

	newUUID, err := uuid.NewV4()
	if err != nil {
		return fmt.Errorf("error pay order in generating uuid: %v", err.Error())
	}

	if err := u.transactionRepo.InsertPayment(tx, &entity.Payment{
		Id:      newUUID.String(),
		OrderId: req.OrderId,
		UserId:  req.UserId,
		Amount:  req.Amount,
		Status:  entity.PaymentStatusPaid,
	}); err != nil {
		return err
	}

	if err := u.transactionRepo.UpdateOrder(tx, &entity.UpdateOrderRequest{
		OrderId: req.OrderId,
		Status:  entity.OrderStatusSucceeded,
	}); err != nil {
		return err
	}

	// post actions execute async (invalidate locks & update stock)
	go func(orderId, userId string) {
		orderItems, err := u.transactionRepo.GetOrderItemsByOrderId(orderId)
		if err != nil {
			log.Printf("error pay order in post action async: error get order item: %v", err.Error())
			return
		}

		for _, item := range orderItems {
			if err := u.redisRepo.InvalidateLockOrderProduct(context.Background(), &entity.LockOrderProductRequest{
				ProductId:   item.ProductId,
				WarehouseId: item.WarehouseId,
				UserId:      userId,
				Quantity:    item.Quantity,
			}); err != nil {
				log.Printf("error pay order in post action async: error invalidate lock product id '%s' for order id '%s': %v", item.ProductId, req.OrderId, err.Error())
			}

			// update product total stock
			if err := u.inventoryRepo.InsertProductWarehouse(&entity.ProductWarehouse{
				ProductId:   item.ProductId,
				WarehouseId: item.WarehouseId,
				TotalStock:  -item.Quantity,
			}); err != nil {
				log.Printf("error pay order in post action async: error update total stock product id '%s' for order id '%s': %v", item.ProductId, req.OrderId, err.Error())
			}
		}
	}(req.OrderId, req.UserId)

	return nil
}
