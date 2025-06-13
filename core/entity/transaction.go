package entity

import (
	"fmt"
	errorutil "mfawzanid/warehouse-commerce/utils/error"
	"time"
)

const (
	OrderStatusPending   = "pending"
	OrderStatusSucceeded = "succeeded"

	OrderExpireTimeInMinute = 1

	PaymentStatusPaid = "paid"
)

type OrderProductItem struct {
	ProductId string `json:"productId"`
	Quantity  int    `json:"quantity"`
}

type OrderProductsRequest struct {
	Items  []*OrderProductItem `json:"items"`
	UserId string
	ShopId string `json:"shopId"`
}

func (r *OrderProductsRequest) Validate() error {
	if len(r.Items) == 0 {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error validate order products request: items are mandatory"))
	}
	return nil
}

type Order struct {
	Id        string
	UserId    string
	ShopId    string
	Amount    int
	Status    string
	CreatedAt time.Time
	ExpiredAt time.Time
}

type LockOrderProductRequest struct {
	ProductId   string
	WarehouseId string
	UserId      string
	Quantity    int
}

type PayOrderRequest struct {
	OrderId string
	Amount  int
	UserId  string
}

type Payment struct {
	Id      string
	OrderId string
	Amount  int
	Status  string
	UserId  string
}

type UpdateOrderRequest struct {
	OrderId string
	Status  string
}

type OrderItem struct {
	OrderId     string
	ProductId   string
	ShopId      string
	WarehouseId string
	Quantity    int
	UnitPrice   int
}

type UpdateProductWarehouseTotalStockRequest struct {
	ProductId   string
	TotalStock  int
	WarehouseId string
}

func (r *UpdateProductWarehouseTotalStockRequest) Validate() error {
	if r.ProductId == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error validate update product stock request: product id is mandatory"))
	}
	if r.WarehouseId == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error validate update product stock request: warehouse id is mandatory"))
	}
	return nil
}

type TransferProductRequest struct {
	ProductId              string `json:"productId"`
	SourceWarehouseId      string `json:"sourceWarehouseId"`
	DestinationWarehouseId string `json:"destinationWarehouseId"`
	TotalStock             int    `json:"totalStock"`
}

func (r *TransferProductRequest) Validate() error {
	if r.ProductId == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error transfer product request: product id is mandatory"))
	}
	if r.SourceWarehouseId == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error transfer product request: source warehouse id is mandatory"))
	}
	if r.DestinationWarehouseId == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error transfer product request: destination warehouse id is mandatory"))
	}
	if r.TotalStock <= 0 {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error transfer product request: total stock must be more than 0"))
	}
	return nil
}
