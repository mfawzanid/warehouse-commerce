package entity

import (
	"errors"
	"fmt"
	errorutil "mfawzanid/warehouse-commerce/utils/error"
)

type CreateWarehouseRequest struct {
	Name string `json:"name"`
}

func (r *CreateWarehouseRequest) Validate() error {
	if r.Name == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error create warehouse request validation: name is mandantory"))
	}
	return nil
}

type Warehouse struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Enabled bool   `json:"enabled"`
}

type UpdateWarehouseStatusRequest struct {
	Id      string
	Enabled bool `json:"enabled"`
}

func (r UpdateWarehouseStatusRequest) Validate() error {
	if r.Id == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error update warehouse status validation: id is mandantory"))
	}
	return nil
}

type GetWarehousesRequest struct {
	Pagination *Pagination
	Enabled    *bool
	Ids        []string
	Name       string
}

func (r GetWarehousesRequest) Validate() error {
	if r.Pagination != nil {
		r.Pagination.Validate()
	}
	return nil
}

type GetWarehousesResponse struct {
	Warehouses []*Warehouse `json:"warehouses"`
	Pagination *Pagination  `json:"pagination"`
}

type CreateShopRequest struct {
	Name string
}

func (r *CreateShopRequest) Validate() error {
	if r.Name == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error create shop request validation: name is mandantory"))
	}
	return nil
}

type Shop struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type GetShopsRequest struct {
	Pagination *Pagination
	Ids        []string
	Name       string
}

func (r *GetShopsRequest) Validate() error {
	if r.Pagination != nil {
		r.Pagination.Validate()
	}
	return nil
}

type GetShopsResponse struct {
	Shops      []*Shop     `json:"shops"`
	Pagination *Pagination `json:"pagination"`
}

type UpsertShopToWarehousesRequest struct {
	ShopId       string   `json:"shopId"`
	WarehouseIds []string `json:"warehousesIds"`
	Enabled      bool     `json:"enabled"`
}

type CreateProductRequest struct {
	Name        string
	Price       int
	TotalStock  int
	WarehouseId string
}

func (r *CreateProductRequest) Validate() error {
	if r.Name == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error create product request validation: name is mandantory"))
	}
	if r.Price <= 0 {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error create product request validation: price must be more than zero"))
	}
	if r.TotalStock <= 0 {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error create product request validation: total stock must be more than zero"))
	}
	if r.WarehouseId == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error create product request validation: name is mandantory"))
	}
	return nil
}

type Product struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	Price int    `json:"price"`
}

type ProductWarehouse struct {
	ProductId   string `json:"productId"`
	TotalStock  int    `json:"totalStock"`
	WarehouseId string `json:"warehouseId"`
}

func (pw ProductWarehouse) Validate() error {
	if pw.ProductId == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, errors.New("error product warehouse validation: product id is mandatory"))
	}
	if pw.WarehouseId == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, errors.New("error product warehouse validation: warehouse id is mandatory"))
	}
	return nil
}

type ProductDetail struct {
	ProductId   string `json:"productId"`
	Name        string `json:"name"`
	Price       int    `json:"price"`
	TotalStock  int    `json:"totalStock"`
	WarehouseId string `json:"warehouseId"`
}

type GetProductWarehousesByQueryRequest struct {
	ProductIds   []string
	WarehouseIds []string
}

func (r *GetProductWarehousesByQueryRequest) Validate() error {
	if len(r.ProductIds) == 0 {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error get product warehouse request validation: product id is mandatory"))
	}
	return nil
}

type GetProductDetailsByShopIdRequest struct {
	ShopId     string
	ProductIds []string
	Pagination *Pagination
}

func (r GetProductDetailsByShopIdRequest) Validate() error {
	if r.ShopId == "" {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, errors.New("error get product details by shop id request validation: shop id is mandatory"))
	}

	// to avoid get all products
	if r.Pagination == nil {
		r.Pagination = &Pagination{}
		r.Pagination.SetToDefault()
	} else {
		r.Pagination.Validate()
	}

	return nil
}

type GetProductDetailsByShopIdResponse struct {
	ProductDetails []*ProductDetail `json:"products"`
	Pagination     *Pagination      `json:"pagination"`
}
