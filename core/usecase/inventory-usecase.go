package usecase

import (
	"fmt"
	"mfawzanid/warehouse-commerce/core/entity"
	"mfawzanid/warehouse-commerce/core/repository"
	errorutil "mfawzanid/warehouse-commerce/utils/error"
	serialutil "mfawzanid/warehouse-commerce/utils/serial"
	transactionutil "mfawzanid/warehouse-commerce/utils/transaction"
)

type InventoryUsecaseInterface interface {
	// warehouse
	CreateWarehouse(req *entity.CreateWarehouseRequest) (string, error)
	UpdateWarehouseStatus(req *entity.UpdateWarehouseStatusRequest) error
	GetWarehouses(req *entity.GetWarehousesRequest) (*entity.GetWarehousesResponse, error)

	// shop
	CreateShop(req *entity.CreateShopRequest) (string, error)
	GetShops(req *entity.GetShopsRequest) (*entity.GetShopsResponse, error)

	// shop-warehouse
	UpsertShopToWarehouses(req *entity.UpsertShopToWarehousesRequest) error

	// product
	CreateProduct(req *entity.CreateProductRequest) (string, error)
	UpdateProductStock(req *entity.UpdateProductWarehouseTotalStockRequest) error
	TransferProduct(req *entity.TransferProductRequest) error
}

type inventoryUsecase struct {
	inventoryRepo repository.InventoryRepositoryInterface
}

func NewInventoryUsecase(inventoryRepo repository.InventoryRepositoryInterface) InventoryUsecaseInterface {
	return &inventoryUsecase{
		inventoryRepo: inventoryRepo,
	}
}

const (
	warehousePrefixSerial = "WRH"
	shopPrefixSerial      = "SHP"
	productPrefixSerial   = "PRD"
)

func (u *inventoryUsecase) CreateWarehouse(req *entity.CreateWarehouseRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}

	warehouseId, err := serialutil.GenerateId(warehousePrefixSerial)
	if err != nil {
		return "", fmt.Errorf("error create warehouse in generating uuid: %v", err.Error())
	}

	warehouse := &entity.Warehouse{
		Id:      warehouseId,
		Name:    req.Name,
		Enabled: false,
	}
	if err := u.inventoryRepo.InsertWarehouse(warehouse); err != nil {
		if err == errorutil.ErrUniqueViolation {
			// warehouse name is unique so if warehouse name has exist, then just return existing warehouse id
			warehouseId, err := u.getWarehouseId(req.Name)
			if err != nil {
				return "", err
			}

			return warehouseId, nil
		} else {
			return "", err
		}
	}

	return warehouse.Id, nil
}

func (u *inventoryUsecase) UpdateWarehouseStatus(req *entity.UpdateWarehouseStatusRequest) error {
	return u.inventoryRepo.UpdateWarehouseStatus(req)
}

func (u *inventoryUsecase) GetWarehouses(req *entity.GetWarehousesRequest) (*entity.GetWarehousesResponse, error) {
	isEnabled := true
	req.Enabled = &isEnabled

	if err := req.Validate(); err != nil {
		return nil, err
	}

	return u.inventoryRepo.GetWarehouses(req)
}

func (u *inventoryUsecase) CreateShop(req *entity.CreateShopRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}

	shopId, err := serialutil.GenerateId(shopPrefixSerial)
	if err != nil {
		return "", fmt.Errorf("error create shop in generating uuid: %v", err.Error())
	}

	shop := &entity.Shop{
		Id:   shopId,
		Name: req.Name,
	}

	if err = u.inventoryRepo.InsertShop(shop); err != nil {
		if err == errorutil.ErrUniqueViolation {
			// shop name is unique so if shop name has exist, then just return existing shop id
			shopId, err := u.getShopId(req.Name)
			if err != nil {
				return "", err
			}

			return shopId, nil
		} else {
			return "", err
		}
	}

	return shop.Id, nil
}

func (u *inventoryUsecase) GetShops(req *entity.GetShopsRequest) (*entity.GetShopsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	return u.inventoryRepo.GetShops(req)
}

func (u *inventoryUsecase) UpsertShopToWarehouses(req *entity.UpsertShopToWarehousesRequest) error {
	// validate shopId whether exist or not
	getShopsResp, err := u.inventoryRepo.GetShops(&entity.GetShopsRequest{
		Ids: []string{req.ShopId},
	})
	if err != nil {
		return err
	}
	if len(getShopsResp.Shops) == 0 {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error upsert shop to warehouses: shop '%v' is not found", req.ShopId))
	}

	// validate warehouseIds whether exist or not
	getWarehousesResp, err := u.inventoryRepo.GetWarehouses(&entity.GetWarehousesRequest{
		Ids: req.WarehouseIds,
	})
	if err != nil {
		return err
	}
	if len(getWarehousesResp.Warehouses) == 0 {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error upsert shop to warehouses: some warehouses are not found"))
	}

	return u.inventoryRepo.InsertShopWarehouses(req)
}

func (u *inventoryUsecase) CreateProduct(req *entity.CreateProductRequest) (string, error) {
	if err := req.Validate(); err != nil {
		return "", err
	}

	// validate warehouseId whether exist or not
	getWarehousesResp, err := u.inventoryRepo.GetWarehouses(&entity.GetWarehousesRequest{
		Ids: []string{req.WarehouseId},
	})
	if err != nil {
		return "", err
	}
	if len(getWarehousesResp.Warehouses) == 0 {
		return "", errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error create product: warehouse '%v' is not found", req.WarehouseId))
	}

	productId, err := serialutil.GenerateId(productPrefixSerial)
	if err != nil {
		return "", fmt.Errorf("error create product in generating uuid: %v", err.Error())
	}

	tx, err := u.inventoryRepo.GetDb().Begin()
	if err != nil {
		return "", fmt.Errorf("error create product in initiating transaction: %v", err.Error())
	}

	defer func() {
		err = transactionutil.SettleTransaction(tx, err)
	}()

	if err := u.inventoryRepo.InsertProduct(tx, &entity.Product{
		Id:    productId,
		Name:  req.Name,
		Price: req.Price,
	}); err != nil {
		if err == errorutil.ErrUniqueViolation {
			// produt name is assumed set as unique so if product name has exist, then just return existing product id
			productId, err = u.getProductId(req.Name)
			if err != nil {
				return "", err
			}

			return productId, nil
		} else {
			return "", err
		}
	}

	if err := u.inventoryRepo.InsertProductWarehouseTx(tx, &entity.ProductWarehouse{
		ProductId:   productId,
		WarehouseId: req.WarehouseId,
		TotalStock:  req.TotalStock,
	}); err != nil {
		return "", err
	}

	return productId, nil
}

func (u *inventoryUsecase) UpdateProductStock(req *entity.UpdateProductWarehouseTotalStockRequest) error {
	return u.inventoryRepo.UpdateProductWarehouseTotalStock(req)
}

func (u *inventoryUsecase) TransferProduct(req *entity.TransferProductRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	// validate both warehouses
	getWarehousesResp, err := u.inventoryRepo.GetWarehouses(&entity.GetWarehousesRequest{
		Ids: []string{req.SourceWarehouseId, req.DestinationWarehouseId},
	})
	if err != nil {
		return err
	}
	if len(getWarehousesResp.Warehouses) != 2 {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error transfer product: some warehouses are not found"))
	}

	// validate product from source warehouse
	sourceProductWarehouses, err := u.inventoryRepo.GetProductWarehousesByQuery(&entity.GetProductWarehousesByQueryRequest{
		ProductIds:   []string{req.ProductId},
		WarehouseIds: []string{req.SourceWarehouseId},
	})
	if err != nil {
		return err
	}
	if len(sourceProductWarehouses) == 0 {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error transfer product: source product warehouse are is not found"))
	}

	sourceProductWarehouse := sourceProductWarehouses[0]
	newSourceTotalStock := sourceProductWarehouse.TotalStock - req.TotalStock
	if newSourceTotalStock < 0 {
		return errorutil.NewErrorCode(errorutil.ErrBadRequest, fmt.Errorf("error transfer product: total stock that transfered is not sufficient"))
	}

	tx, err := u.inventoryRepo.GetDb().Begin()
	if err != nil {
		return fmt.Errorf("error transfer product in initiating transaction: %v", err.Error())
	}

	defer func() {
		err = transactionutil.SettleTransaction(tx, err)
	}()

	if err := u.inventoryRepo.InsertProductWarehouseTx(tx, &entity.ProductWarehouse{
		ProductId:   req.ProductId,
		WarehouseId: req.SourceWarehouseId,
		TotalStock:  -req.TotalStock, // add existing value with -totalStock that will transfer
	}); err != nil {
		return err
	}

	if err := u.inventoryRepo.InsertProductWarehouseTx(tx, &entity.ProductWarehouse{
		ProductId:   req.ProductId,
		WarehouseId: req.DestinationWarehouseId,
		TotalStock:  req.TotalStock, // add existing value with totalStock that will transfer
	}); err != nil {
		return err
	}

	return nil
}
