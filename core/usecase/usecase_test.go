package usecase_test

import (
	"errors"
	"mfawzanid/warehouse-commerce/core/entity"
	"mfawzanid/warehouse-commerce/core/mocks"
	"mfawzanid/warehouse-commerce/core/usecase"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type usecaseTest struct {
	userRepo        *mocks.UserRepositoryInterface
	inventoryRepo   *mocks.InventoryRepositoryInterface
	transactionRepo *mocks.TransactionRepositoryInterface
	redisRepo       *mocks.RedisRepositoryInterface

	authUsecase        usecase.AuthUsecaseInterface
	userUsecase        usecase.UserUsecaseInterface
	inventoryUsecase   usecase.InventoryUsecaseInterface
	transactionUsecase usecase.TransactionUsecaseInterface
}

var ucTest usecaseTest

func init() {
	mockUserRepo := mocks.UserRepositoryInterface{}
	mockInventoryRepo := mocks.InventoryRepositoryInterface{}
	mockTransactionRepo := mocks.TransactionRepositoryInterface{}
	mockRedisRepo := mocks.RedisRepositoryInterface{}

	authUsecase := usecase.NewAuthUsecase()
	userUsecase := usecase.NewUserUsecase(&mockUserRepo, authUsecase)
	inventoryUsecase := usecase.NewInventoryUsecase(&mockInventoryRepo)
	transactionUsecase := usecase.NewTransactionUsecase(&mockInventoryRepo, &mockTransactionRepo, &mockRedisRepo)

	ucTest = usecaseTest{
		userRepo:        &mockUserRepo,
		inventoryRepo:   &mockInventoryRepo,
		transactionRepo: &mockTransactionRepo,
		redisRepo:       &mockRedisRepo,

		authUsecase:        authUsecase,
		userUsecase:        userUsecase,
		inventoryUsecase:   inventoryUsecase,
		transactionUsecase: transactionUsecase,
	}
}

func TestRegisterUser(t *testing.T) {
	t.Run("RegisterUser_identifier type not valid_then return error", func(t *testing.T) {
		email := "email_1@mail.com"

		token, err := ucTest.userUsecase.RegisterUser(&entity.RegisterUserRequest{
			IdentifierType: "xxx",
			Identifier:     email,
		})

		assert.NotNil(t, err)
		assert.Empty(t, token)
	})
	t.Run("RegisterUser_identifier value is empty_then return error", func(t *testing.T) {
		token, err := ucTest.userUsecase.RegisterUser(&entity.RegisterUserRequest{
			IdentifierType: "email",
			Identifier:     "",
		})

		assert.NotNil(t, err)
		assert.Empty(t, token)
	})
	t.Run("RegisterUser_correct payload_then return success", func(t *testing.T) {
		email := "email_1@mail.com"
		ucTest.userRepo.On("InsertUser", mock.Anything).Return(nil).Once()

		token, err := ucTest.userUsecase.RegisterUser(&entity.RegisterUserRequest{
			IdentifierType: "email",
			Identifier:     email,
		})

		assert.Nil(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestLogin(t *testing.T) {
	t.Run("Login_identifier type not valid_then return error", func(t *testing.T) {
		email := "email_1@mail.com"

		token, err := ucTest.userUsecase.Login(&entity.LoginRequest{
			IdentifierType: "xxx",
			Identifier:     email,
		})

		assert.NotNil(t, err)
		assert.Empty(t, token)
	})
	t.Run("Login_identifier value is empty_then return error", func(t *testing.T) {
		token, err := ucTest.userUsecase.Login(&entity.LoginRequest{
			IdentifierType: "email",
			Identifier:     "",
		})

		assert.NotNil(t, err)
		assert.Empty(t, token)
	})
	t.Run("Login_correct payload_then return success", func(t *testing.T) {
		email := "email_1@mail.com"
		expectedUser := &entity.User{
			Email: email,
		}

		ucTest.userRepo.On("GetUser", mock.Anything).Return(expectedUser, nil).Once()

		token, err := ucTest.userUsecase.Login(&entity.LoginRequest{
			IdentifierType: "email",
			Identifier:     email,
		})

		assert.Nil(t, err)
		assert.NotEmpty(t, token)
	})
}

func TestCreateWarehouse(t *testing.T) {
	t.Run("CreateWarehouse_empty warehouse name_then return success", func(t *testing.T) {
		name := ""

		id, err := ucTest.inventoryUsecase.CreateWarehouse(&entity.CreateWarehouseRequest{
			Name: name,
		})

		assert.NotNil(t, err)
		assert.Empty(t, id)
	})
	t.Run("CreateWarehouse_correct payload_then return success", func(t *testing.T) {
		name := "name"

		ucTest.inventoryRepo.On("InsertWarehouse", mock.Anything).Return(nil).Once()

		id, err := ucTest.inventoryUsecase.CreateWarehouse(&entity.CreateWarehouseRequest{
			Name: name,
		})

		assert.Nil(t, err)
		assert.NotEmpty(t, id)
	})
}

func TestUpdateWarehouseStatus(t *testing.T) {
	t.Run("UpdateWarehouseStatus_correct payload_then return success", func(t *testing.T) {
		id := "id"

		ucTest.inventoryRepo.On("UpdateWarehouseStatus", mock.Anything).Return(nil).Once()

		err := ucTest.inventoryUsecase.UpdateWarehouseStatus(&entity.UpdateWarehouseStatusRequest{
			Id: id,
		})

		assert.Nil(t, err)
	})
}

func TestGetWarehouses(t *testing.T) {
	t.Run("GetWarehouses_correct payload_then return success", func(t *testing.T) {
		ids := []string{"id_1"}
		ucTest.inventoryRepo.On("GetWarehouses", mock.Anything).Return(&entity.GetWarehousesResponse{
			Warehouses: []*entity.Warehouse{},
			Pagination: &entity.Pagination{},
		}, nil).Once()

		resp, err := ucTest.inventoryUsecase.GetWarehouses(&entity.GetWarehousesRequest{
			Ids: ids,
		})

		assert.Nil(t, err)
		assert.NotEmpty(t, resp)
	})
}

func TestCreateShop(t *testing.T) {
	t.Run("CreateShop_empty name_then return error", func(t *testing.T) {
		name := ""

		shopId, err := ucTest.inventoryUsecase.CreateShop(&entity.CreateShopRequest{
			Name: name,
		})

		assert.NotNil(t, err)
		assert.Empty(t, shopId)
	})
	t.Run("CreateShop_correct payload_then return success", func(t *testing.T) {
		name := "shop_1"

		ucTest.inventoryRepo.On("InsertShop", mock.Anything).Return(nil).Once()

		shopId, err := ucTest.inventoryUsecase.CreateShop(&entity.CreateShopRequest{
			Name: name,
		})

		assert.Nil(t, err)
		assert.NotEmpty(t, shopId)
	})
}

func TestGetShops(t *testing.T) {
	t.Run("GetShops_correct payload_then return success", func(t *testing.T) {
		ids := []string{"id_1"}

		ucTest.inventoryRepo.On("GetShops", mock.Anything).Return(&entity.GetShopsResponse{
			Shops:      []*entity.Shop{},
			Pagination: &entity.Pagination{},
		}, nil).Once()

		resp, err := ucTest.inventoryUsecase.GetShops(&entity.GetShopsRequest{
			Ids: ids,
		})

		assert.Nil(t, err)
		assert.NotEmpty(t, resp)
	})
}

func TestUpsertShopToWarehouses(t *testing.T) {
	t.Run("UpsertShopToWarehouses_shop id is not found_then return error", func(t *testing.T) {
		ucTest.inventoryRepo.On("GetShops", mock.Anything).Return(nil, errors.New("")).Once()

		err := ucTest.inventoryUsecase.UpsertShopToWarehouses(&entity.UpsertShopToWarehousesRequest{})

		assert.NotNil(t, err)
	})
	t.Run("UpsertShopToWarehouses_some warehouse id is not found_then return error", func(t *testing.T) {
		shops := []*entity.Shop{{Id: "id"}}

		ucTest.inventoryRepo.On("GetShops", mock.Anything).Return(&entity.GetShopsResponse{
			Shops: shops,
		}, nil).Once()

		ucTest.inventoryRepo.On("GetWarehouses", mock.Anything).Return(nil, errors.New("")).Once()

		err := ucTest.inventoryUsecase.UpsertShopToWarehouses(&entity.UpsertShopToWarehousesRequest{})

		assert.NotNil(t, err)
	})

	t.Run("UpsertShopToWarehouses_correct payload_then successfully updated", func(t *testing.T) {
		shops := []*entity.Shop{{Id: "id"}}
		ucTest.inventoryRepo.On("GetShops", mock.Anything).Return(&entity.GetShopsResponse{
			Shops: shops,
		}, nil).Once()

		warehouses := []*entity.Warehouse{{Id: "id"}}
		ucTest.inventoryRepo.On("GetWarehouses", mock.Anything).Return(&entity.GetWarehousesResponse{
			Warehouses: warehouses,
		}, nil).Once()

		ucTest.inventoryRepo.On("InsertShopWarehouses", mock.Anything).Return(nil).Once()

		err := ucTest.inventoryUsecase.UpsertShopToWarehouses(&entity.UpsertShopToWarehousesRequest{})

		assert.Nil(t, err)
	})
}

func TestCreateProduct(t *testing.T) {
	t.Run("CreateProduct_bad request_then return error", func(t *testing.T) {
		req := &entity.CreateProductRequest{}

		_, err := ucTest.inventoryUsecase.CreateProduct(req)

		assert.NotNil(t, err)
	})
	t.Run("CreateProduct_get warehouse is error_then return error", func(t *testing.T) {
		warehouseId := "warehouse_id"
		req := &entity.CreateProductRequest{
			Name:        "name",
			Price:       1000,
			TotalStock:  10,
			WarehouseId: warehouseId,
		}

		warehouses := []*entity.Warehouse{{Id: warehouseId}}
		ucTest.inventoryRepo.On("GetWarehouses", mock.Anything).Return(&entity.GetWarehousesResponse{
			Warehouses: warehouses,
		}, errors.New("")).Once()

		productId, err := ucTest.inventoryUsecase.CreateProduct(req)

		assert.NotNil(t, err)
		assert.Empty(t, productId)
	})
	t.Run("CreateProduct_insert product is error_then return error", func(t *testing.T) {
		// setup sqlmock for transaction
		db, mockDB, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		warehouseId := "warehouse_id"
		req := &entity.CreateProductRequest{
			Name:        "name",
			Price:       1000,
			TotalStock:  10,
			WarehouseId: warehouseId,
		}

		warehouses := []*entity.Warehouse{{Id: warehouseId}}
		ucTest.inventoryRepo.On("GetWarehouses", mock.Anything).Return(&entity.GetWarehousesResponse{
			Warehouses: warehouses,
		}, nil).Once()

		ucTest.inventoryRepo.On("GetDb").Return(db).Once()
		mockDB.ExpectBegin()

		ucTest.inventoryRepo.On("InsertProduct", mock.Anything, mock.AnythingOfType("*entity.Product")).Return(errors.New("")).Once()

		productId, err := ucTest.inventoryUsecase.CreateProduct(req)

		assert.NotNil(t, err)
		assert.Empty(t, productId)
	})
	t.Run("CreateProduct_insert product warehouse is error_then return error", func(t *testing.T) {
		// setup sqlmock for transaction
		db, mockDB, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		warehouseId := "warehouse_id"
		req := &entity.CreateProductRequest{
			Name:        "name",
			Price:       1000,
			TotalStock:  10,
			WarehouseId: warehouseId,
		}

		warehouses := []*entity.Warehouse{{Id: warehouseId}}
		ucTest.inventoryRepo.On("GetWarehouses", mock.Anything).Return(&entity.GetWarehousesResponse{
			Warehouses: warehouses,
		}, nil).Once()

		ucTest.inventoryRepo.On("GetDb").Return(db).Once()
		mockDB.ExpectBegin()

		ucTest.inventoryRepo.On("InsertProduct", mock.Anything, mock.AnythingOfType("*entity.Product")).Return(nil).Once()
		ucTest.inventoryRepo.On("InsertProductWarehouseTx", mock.Anything, mock.AnythingOfType("*entity.ProductWarehouse")).Return(errors.New("")).Once()

		productId, err := ucTest.inventoryUsecase.CreateProduct(req)

		assert.NotNil(t, err)
		assert.Empty(t, productId)
	})
	t.Run("CreateProduct_correct payload_then return success", func(t *testing.T) {
		// setup sqlmock for transaction
		db, mockDB, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		warehouseId := "warehouse_id"
		req := &entity.CreateProductRequest{
			Name:        "name",
			Price:       1000,
			TotalStock:  10,
			WarehouseId: warehouseId,
		}

		warehouses := []*entity.Warehouse{{Id: warehouseId}}
		ucTest.inventoryRepo.On("GetWarehouses", mock.Anything).Return(&entity.GetWarehousesResponse{
			Warehouses: warehouses,
		}, nil).Once()

		ucTest.inventoryRepo.On("GetDb").Return(db).Once()
		mockDB.ExpectBegin()

		ucTest.inventoryRepo.On("InsertProduct", mock.Anything, mock.AnythingOfType("*entity.Product")).Return(nil).Once()
		ucTest.inventoryRepo.On("InsertProductWarehouseTx", mock.Anything, mock.AnythingOfType("*entity.ProductWarehouse")).Return(nil).Once()

		productId, err := ucTest.inventoryUsecase.CreateProduct(req)

		assert.NoError(t, err)
		assert.NotEmpty(t, productId)
	})
}

func TestUpdateProductStock(t *testing.T) {
	t.Run("UpdateProductStock_correct payload_then return success", func(t *testing.T) {
		req := &entity.UpdateProductWarehouseTotalStockRequest{
			ProductId:   "product_id",
			WarehouseId: "warehouse_id",
			TotalStock:  10,
		}

		ucTest.inventoryRepo.On("UpdateProductWarehouseTotalStock", req).Return(nil).Once()

		err := ucTest.inventoryUsecase.UpdateProductStock(req)

		assert.Nil(t, err)
	})
}

func TestTransferProduct(t *testing.T) {
	t.Run("TransferProduct_bad request_then return error", func(t *testing.T) {
		req := &entity.TransferProductRequest{}

		err := ucTest.inventoryUsecase.TransferProduct(req)

		assert.NotNil(t, err)
	})
	t.Run("TransferProduct_get warehouses is error_then return error", func(t *testing.T) {
		productId := "product_id"
		sourceWarehouseId := "source_warehouse_id"
		destinationWarehouseId := "destination_warehouse_id"

		req := &entity.TransferProductRequest{
			ProductId:              productId,
			SourceWarehouseId:      sourceWarehouseId,
			DestinationWarehouseId: destinationWarehouseId,
			TotalStock:             10,
		}

		// mock GetWarehouses
		warehouses := []*entity.Warehouse{
			{Id: sourceWarehouseId},
			{Id: destinationWarehouseId},
		}
		ucTest.inventoryRepo.On("GetWarehouses", &entity.GetWarehousesRequest{Ids: []string{sourceWarehouseId, destinationWarehouseId}}).Return(&entity.GetWarehousesResponse{
			Warehouses: warehouses,
		}, errors.New("")).Once()

		err := ucTest.inventoryUsecase.TransferProduct(req)

		assert.NotNil(t, err)
	})
	t.Run("TransferProduct_get product warehouses is error_then return error", func(t *testing.T) {
		productId := "product_id"
		sourceWarehouseId := "source_warehouse_id"
		destinationWarehouseId := "destination_warehouse_id"

		req := &entity.TransferProductRequest{
			ProductId:              productId,
			SourceWarehouseId:      sourceWarehouseId,
			DestinationWarehouseId: destinationWarehouseId,
			TotalStock:             10,
		}

		// mock GetWarehouses
		warehouses := []*entity.Warehouse{
			{Id: sourceWarehouseId},
			{Id: destinationWarehouseId},
		}
		ucTest.inventoryRepo.On("GetWarehouses", &entity.GetWarehousesRequest{Ids: []string{sourceWarehouseId, destinationWarehouseId}}).Return(&entity.GetWarehousesResponse{
			Warehouses: warehouses,
		}, nil).Once()

		ucTest.inventoryRepo.On("GetProductWarehousesByQuery", &entity.GetProductWarehousesByQueryRequest{
			ProductIds:   []string{productId},
			WarehouseIds: []string{sourceWarehouseId},
		}).Return(nil, errors.New("")).Once()

		err := ucTest.inventoryUsecase.TransferProduct(req)

		assert.NotNil(t, err)
	})
	t.Run("TransferProduct_total stock from source is not sufficient_then return error", func(t *testing.T) {
		productId := "product_id"
		sourceWarehouseId := "source_warehouse_id"
		destinationWarehouseId := "destination_warehouse_id"

		req := &entity.TransferProductRequest{
			ProductId:              productId,
			SourceWarehouseId:      sourceWarehouseId,
			DestinationWarehouseId: destinationWarehouseId,
			TotalStock:             10,
		}

		// mock GetWarehouses
		warehouses := []*entity.Warehouse{
			{Id: sourceWarehouseId},
			{Id: destinationWarehouseId},
		}
		ucTest.inventoryRepo.On("GetWarehouses", &entity.GetWarehousesRequest{Ids: []string{sourceWarehouseId, destinationWarehouseId}}).Return(&entity.GetWarehousesResponse{
			Warehouses: warehouses,
		}, nil).Once()

		// mock GetProductWarehousesByQuery
		sourceProductWarehouse := &entity.ProductWarehouse{
			ProductId:  productId,
			TotalStock: 1,
		}
		productWarehouses := []*entity.ProductWarehouse{sourceProductWarehouse}
		ucTest.inventoryRepo.On("GetProductWarehousesByQuery", &entity.GetProductWarehousesByQueryRequest{
			ProductIds:   []string{productId},
			WarehouseIds: []string{sourceWarehouseId},
		}).Return(productWarehouses, nil).Once()

		err := ucTest.inventoryUsecase.TransferProduct(req)

		assert.NotNil(t, err)
	})
	t.Run("TransferProduct_insert product to source warehouse error_then return error", func(t *testing.T) {
		productId := "product_id"
		sourceWarehouseId := "source_warehouse_id"
		destinationWarehouseId := "destination_warehouse_id"

		req := &entity.TransferProductRequest{
			ProductId:              productId,
			SourceWarehouseId:      sourceWarehouseId,
			DestinationWarehouseId: destinationWarehouseId,
			TotalStock:             10,
		}

		// mock GetWarehouses
		warehouses := []*entity.Warehouse{
			{Id: sourceWarehouseId},
			{Id: destinationWarehouseId},
		}
		ucTest.inventoryRepo.On("GetWarehouses", &entity.GetWarehousesRequest{Ids: []string{sourceWarehouseId, destinationWarehouseId}}).Return(&entity.GetWarehousesResponse{
			Warehouses: warehouses,
		}, nil).Once()

		// mock GetProductWarehousesByQuery
		sourceProductWarehouse := &entity.ProductWarehouse{
			ProductId:  productId,
			TotalStock: 100,
		}
		productWarehouses := []*entity.ProductWarehouse{sourceProductWarehouse}
		ucTest.inventoryRepo.On("GetProductWarehousesByQuery", &entity.GetProductWarehousesByQueryRequest{
			ProductIds:   []string{productId},
			WarehouseIds: []string{sourceWarehouseId},
		}).Return(productWarehouses, nil).Once()

		// setup sqlmock for transaction
		db, mockDB, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ucTest.inventoryRepo.On("GetDb").Return(db).Once()
		mockDB.ExpectBegin()

		ucTest.inventoryRepo.On("InsertProductWarehouseTx", mock.Anything, mock.Anything).Return(errors.New("")).Once()

		err = ucTest.inventoryUsecase.TransferProduct(req)

		assert.NotNil(t, err)
	})
	t.Run("TransferProduct_insert product warehouse error_then return error", func(t *testing.T) {
		productId := "product_id"
		sourceWarehouseId := "source_warehouse_id"
		destinationWarehouseId := "destination_warehouse_id"

		req := &entity.TransferProductRequest{
			ProductId:              productId,
			SourceWarehouseId:      sourceWarehouseId,
			DestinationWarehouseId: destinationWarehouseId,
			TotalStock:             10,
		}

		// mock GetWarehouses
		warehouses := []*entity.Warehouse{
			{Id: sourceWarehouseId},
			{Id: destinationWarehouseId},
		}
		ucTest.inventoryRepo.On("GetWarehouses", &entity.GetWarehousesRequest{Ids: []string{sourceWarehouseId, destinationWarehouseId}}).Return(&entity.GetWarehousesResponse{
			Warehouses: warehouses,
		}, nil).Once()

		// mock GetProductWarehousesByQuery
		sourceProductWarehouse := &entity.ProductWarehouse{
			ProductId:  productId,
			TotalStock: 100,
		}
		productWarehouses := []*entity.ProductWarehouse{sourceProductWarehouse}
		ucTest.inventoryRepo.On("GetProductWarehousesByQuery", &entity.GetProductWarehousesByQueryRequest{
			ProductIds:   []string{productId},
			WarehouseIds: []string{sourceWarehouseId},
		}).Return(productWarehouses, nil).Once()

		// setup sqlmock for transaction
		db, mockDB, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ucTest.inventoryRepo.On("GetDb").Return(db).Once()
		mockDB.ExpectBegin()

		ucTest.inventoryRepo.On("InsertProductWarehouseTx", mock.Anything, mock.Anything).Return(nil).Once()
		ucTest.inventoryRepo.On("InsertProductWarehouseTx", mock.Anything, mock.Anything).Return(errors.New("")).Once()

		err = ucTest.inventoryUsecase.TransferProduct(req)

		assert.NotNil(t, err)
	})

	t.Run("TransferProduct_correct payload_then return success", func(t *testing.T) {
		productId := "product_id"
		sourceWarehouseId := "source_warehouse_id"
		destinationWarehouseId := "destination_warehouse_id"

		req := &entity.TransferProductRequest{
			ProductId:              productId,
			SourceWarehouseId:      sourceWarehouseId,
			DestinationWarehouseId: destinationWarehouseId,
			TotalStock:             10,
		}

		// mock GetWarehouses
		warehouses := []*entity.Warehouse{
			{Id: sourceWarehouseId},
			{Id: destinationWarehouseId},
		}
		ucTest.inventoryRepo.On("GetWarehouses", &entity.GetWarehousesRequest{Ids: []string{sourceWarehouseId, destinationWarehouseId}}).Return(&entity.GetWarehousesResponse{
			Warehouses: warehouses,
		}, nil).Once()

		// mock GetProductWarehousesByQuery
		sourceProductWarehouse := &entity.ProductWarehouse{
			ProductId:  productId,
			TotalStock: 100,
		}
		productWarehouses := []*entity.ProductWarehouse{sourceProductWarehouse}
		ucTest.inventoryRepo.On("GetProductWarehousesByQuery", &entity.GetProductWarehousesByQueryRequest{
			ProductIds:   []string{productId},
			WarehouseIds: []string{sourceWarehouseId},
		}).Return(productWarehouses, nil).Once()

		// setup sqlmock for transaction
		db, mockDB, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ucTest.inventoryRepo.On("GetDb").Return(db).Once()
		mockDB.ExpectBegin()

		ucTest.inventoryRepo.On("InsertProductWarehouseTx", mock.Anything, mock.Anything).Return(nil).Once()
		ucTest.inventoryRepo.On("InsertProductWarehouseTx", mock.Anything, mock.Anything).Return(nil).Once()

		err = ucTest.inventoryUsecase.TransferProduct(req)

		assert.NoError(t, err)
	})
}

func TestGetProductDetailsByShopId(t *testing.T) {
	t.Run("GetProductDetailsByShopId_bad request_then return error", func(t *testing.T) {
		resp, err := ucTest.transactionUsecase.GetProductDetailsByShopId(&entity.GetProductDetailsByShopIdRequest{})

		assert.NotNil(t, err)
		assert.Empty(t, resp)
	})
	t.Run("GetProductDetailsByShopId_get product details error_then return error", func(t *testing.T) {
		shopId := "shopId"

		ucTest.inventoryRepo.On("GetProductDetailsByShopId", mock.Anything).Return(nil, errors.New("")).Once()

		resp, err := ucTest.transactionUsecase.GetProductDetailsByShopId(&entity.GetProductDetailsByShopIdRequest{
			ShopId: shopId,
		})

		assert.NotNil(t, err)
		assert.Empty(t, resp)
	})
	t.Run("GetProductDetailsByShopId_get reserved product error_then return error", func(t *testing.T) {
		shopId := "shopId"
		productId := "productId"
		warehouseId := "warehouseId"
		totalStock := 100

		productDetails := []*entity.ProductDetail{
			{
				ProductId:   productId,
				WarehouseId: warehouseId,
				TotalStock:  totalStock,
			},
		}
		getProductResp := &entity.GetProductDetailsByShopIdResponse{
			ProductDetails: productDetails,
		}

		ucTest.inventoryRepo.On("GetProductDetailsByShopId", mock.Anything).Return(getProductResp, nil).Once()
		ucTest.redisRepo.On("GetReservedProductQuantity", mock.Anything, productId, warehouseId).Return(0, errors.New("")).Once()

		resp, err := ucTest.transactionUsecase.GetProductDetailsByShopId(&entity.GetProductDetailsByShopIdRequest{
			ShopId: shopId,
		})

		assert.NotNil(t, err)
		assert.Empty(t, resp)
	})
	t.Run("GetProductDetailsByShopId_correct payload_then return success", func(t *testing.T) {
		shopId := "shopId"
		productId := "productId"
		warehouseId := "warehouseId"
		totalStock := 100
		reservedStock := 5

		productDetails := []*entity.ProductDetail{
			{
				ProductId:   productId,
				WarehouseId: warehouseId,
				TotalStock:  totalStock,
			},
		}
		getProductResp := &entity.GetProductDetailsByShopIdResponse{
			ProductDetails: productDetails,
		}

		expectedProductDetails := []*entity.ProductDetail{
			{
				ProductId:   productId,
				WarehouseId: warehouseId,
				TotalStock:  totalStock - reservedStock,
			},
		}

		ucTest.inventoryRepo.On("GetProductDetailsByShopId", mock.Anything).Return(getProductResp, nil).Once()
		ucTest.redisRepo.On("GetReservedProductQuantity", mock.Anything, productId, warehouseId).Return(reservedStock, nil).Once()

		resp, err := ucTest.transactionUsecase.GetProductDetailsByShopId(&entity.GetProductDetailsByShopIdRequest{
			ShopId: shopId,
		})

		assert.Nil(t, err)
		assert.NotEmpty(t, resp)
		assert.EqualValues(t, expectedProductDetails, resp.ProductDetails)
	})
}

func TestOrderProducts(t *testing.T) {
	t.Run("OrderProducts_bad request_then return error", func(t *testing.T) {
		req := &entity.OrderProductsRequest{}
		id, err := ucTest.transactionUsecase.OrderProducts(req)

		assert.NotNil(t, err)
		assert.Empty(t, id)
	})
	t.Run("OrderProducts_get product detail by shop id error_then return error", func(t *testing.T) {
		productId := "productId"
		shopId := "shopId"

		// mock GetProductDetailsByShopId
		ucTest.inventoryRepo.On("GetProductDetailsByShopId", &entity.GetProductDetailsByShopIdRequest{
			ShopId:     shopId,
			ProductIds: []string{productId},
		}).Return(nil, errors.New("")).Once()

		// usecase
		orderProductItem := &entity.OrderProductItem{
			ProductId: productId,
			Quantity:  5,
		}
		req := &entity.OrderProductsRequest{
			Items:  []*entity.OrderProductItem{orderProductItem},
			ShopId: shopId,
			UserId: "userId",
		}
		id, err := ucTest.transactionUsecase.OrderProducts(req)

		assert.NotNil(t, err)
		assert.Empty(t, id)
	})
	t.Run("OrderProducts_get reserved product quantity error_then return error", func(t *testing.T) {
		productId := "productId"
		shopId := "shopId"
		warehouseId := "warehouseId"

		// mock GetProductDetailsByShopId
		getProductDetailsByShopIdResp := &entity.GetProductDetailsByShopIdResponse{
			ProductDetails: []*entity.ProductDetail{
				{
					ProductId:   productId,
					WarehouseId: warehouseId,
					TotalStock:  100,
				},
			},
		}

		ucTest.inventoryRepo.On("GetProductDetailsByShopId", &entity.GetProductDetailsByShopIdRequest{
			ShopId:     shopId,
			ProductIds: []string{productId},
		}).Return(getProductDetailsByShopIdResp, nil).Once()

		// mock GetReservedProductQuantity
		ucTest.redisRepo.On("GetReservedProductQuantity", mock.Anything, productId, warehouseId).Return(0, errors.New("")).Once()

		// usecase
		orderProductItem := &entity.OrderProductItem{
			ProductId: productId,
			Quantity:  5,
		}
		req := &entity.OrderProductsRequest{
			Items:  []*entity.OrderProductItem{orderProductItem},
			ShopId: shopId,
			UserId: "userId",
		}
		id, err := ucTest.transactionUsecase.OrderProducts(req)

		assert.NotNil(t, err)
		assert.Empty(t, id)
	})
	t.Run("OrderProducts_lock order product error_then return error", func(t *testing.T) {
		productId := "productId"
		shopId := "shopId"
		warehouseId := "warehouseId"

		// mock GetProductDetailsByShopId
		getProductDetailsByShopIdResp := &entity.GetProductDetailsByShopIdResponse{
			ProductDetails: []*entity.ProductDetail{
				{
					ProductId:   productId,
					WarehouseId: warehouseId,
					TotalStock:  100,
				},
			},
		}

		ucTest.inventoryRepo.On("GetProductDetailsByShopId", &entity.GetProductDetailsByShopIdRequest{
			ShopId:     shopId,
			ProductIds: []string{productId},
		}).Return(getProductDetailsByShopIdResp, nil).Once()

		// mock GetReservedProductQuantity
		reservedQuantity := 5
		ucTest.redisRepo.On("GetReservedProductQuantity", mock.Anything, productId, warehouseId).Return(reservedQuantity, nil).Once()

		// mock LockOrderProduct
		ucTest.redisRepo.On("LockOrderProduct", mock.Anything, mock.Anything).Return(errors.New("")).Once()

		// usecase
		orderProductItem := &entity.OrderProductItem{
			ProductId: productId,
			Quantity:  5,
		}
		req := &entity.OrderProductsRequest{
			Items:  []*entity.OrderProductItem{orderProductItem},
			ShopId: shopId,
			UserId: "userId",
		}
		id, err := ucTest.transactionUsecase.OrderProducts(req)

		assert.NotNil(t, err)
		assert.Empty(t, id)
	})
	t.Run("OrderProducts_insert order error_then return error", func(t *testing.T) {
		productId := "productId"
		shopId := "shopId"
		warehouseId := "warehouseId"

		// mock GetProductDetailsByShopId
		getProductDetailsByShopIdResp := &entity.GetProductDetailsByShopIdResponse{
			ProductDetails: []*entity.ProductDetail{
				{
					ProductId:   productId,
					WarehouseId: warehouseId,
					TotalStock:  100,
				},
			},
		}

		ucTest.inventoryRepo.On("GetProductDetailsByShopId", &entity.GetProductDetailsByShopIdRequest{
			ShopId:     shopId,
			ProductIds: []string{productId},
		}).Return(getProductDetailsByShopIdResp, nil).Once()

		// mock GetReservedProductQuantity
		reservedQuantity := 5
		ucTest.redisRepo.On("GetReservedProductQuantity", mock.Anything, productId, warehouseId).Return(reservedQuantity, nil).Once()

		// mock LockOrderProduct
		ucTest.redisRepo.On("LockOrderProduct", mock.Anything, mock.Anything).Return(nil).Once()

		// mock InsertOrder
		ucTest.transactionRepo.On("InsertOrder", mock.Anything).Return(errors.New("")).Once()

		// usecase
		orderProductItem := &entity.OrderProductItem{
			ProductId: productId,
			Quantity:  5,
		}
		req := &entity.OrderProductsRequest{
			Items:  []*entity.OrderProductItem{orderProductItem},
			ShopId: shopId,
			UserId: "userId",
		}
		id, err := ucTest.transactionUsecase.OrderProducts(req)

		assert.NotNil(t, err)
		assert.Empty(t, id)
	})
	t.Run("OrderProducts_insert order item error_then return error", func(t *testing.T) {
		productId := "productId"
		shopId := "shopId"
		warehouseId := "warehouseId"

		// mock GetProductDetailsByShopId
		getProductDetailsByShopIdResp := &entity.GetProductDetailsByShopIdResponse{
			ProductDetails: []*entity.ProductDetail{
				{
					ProductId:   productId,
					WarehouseId: warehouseId,
					TotalStock:  100,
				},
			},
		}

		ucTest.inventoryRepo.On("GetProductDetailsByShopId", &entity.GetProductDetailsByShopIdRequest{
			ShopId:     shopId,
			ProductIds: []string{productId},
		}).Return(getProductDetailsByShopIdResp, nil).Once()

		// mock GetReservedProductQuantity
		reservedQuantity := 5
		ucTest.redisRepo.On("GetReservedProductQuantity", mock.Anything, productId, warehouseId).Return(reservedQuantity, nil).Once()

		// mock LockOrderProduct
		ucTest.redisRepo.On("LockOrderProduct", mock.Anything, mock.Anything).Return(nil).Once()

		// mock InsertOrder
		ucTest.transactionRepo.On("InsertOrder", mock.Anything).Return(nil).Once()

		// mock InsertOrderItems
		ucTest.transactionRepo.On("InsertOrderItems", mock.Anything).Return(errors.New("")).Once()

		// usecase
		orderProductItem := &entity.OrderProductItem{
			ProductId: productId,
			Quantity:  5,
		}
		req := &entity.OrderProductsRequest{
			Items:  []*entity.OrderProductItem{orderProductItem},
			ShopId: shopId,
			UserId: "userId",
		}
		id, err := ucTest.transactionUsecase.OrderProducts(req)

		assert.NotNil(t, err)
		assert.Empty(t, id)
	})
	t.Run("OrderProducts_correct payload_then return success", func(t *testing.T) {
		productId := "productId"
		shopId := "shopId"
		warehouseId := "warehouseId"

		// mock GetProductDetailsByShopId
		getProductDetailsByShopIdResp := &entity.GetProductDetailsByShopIdResponse{
			ProductDetails: []*entity.ProductDetail{
				{
					ProductId:   productId,
					WarehouseId: warehouseId,
					TotalStock:  100,
				},
			},
		}

		ucTest.inventoryRepo.On("GetProductDetailsByShopId", &entity.GetProductDetailsByShopIdRequest{
			ShopId:     shopId,
			ProductIds: []string{productId},
		}).Return(getProductDetailsByShopIdResp, nil).Once()

		// mock GetReservedProductQuantity
		reservedQuantity := 5
		ucTest.redisRepo.On("GetReservedProductQuantity", mock.Anything, productId, warehouseId).Return(reservedQuantity, nil).Once()

		// mock LockOrderProduct
		ucTest.redisRepo.On("LockOrderProduct", mock.Anything, mock.Anything).Return(nil).Once()

		// mock InsertOrder
		ucTest.transactionRepo.On("InsertOrder", mock.Anything).Return(nil).Once()

		// mock InsertOrderItems
		ucTest.transactionRepo.On("InsertOrderItems", mock.Anything).Return(nil).Once()

		// usecase
		orderProductItem := &entity.OrderProductItem{
			ProductId: productId,
			Quantity:  5,
		}
		req := &entity.OrderProductsRequest{
			Items:  []*entity.OrderProductItem{orderProductItem},
			ShopId: shopId,
			UserId: "userId",
		}
		id, err := ucTest.transactionUsecase.OrderProducts(req)

		assert.Nil(t, err)
		assert.NotEmpty(t, id)
	})
}

func TestPayOrder(t *testing.T) {
	t.Run("PayOrder_bad request_then return error", func(t *testing.T) {
		// usecase
		err := ucTest.transactionUsecase.PayOrder(&entity.PayOrderRequest{})

		assert.NotNil(t, err)
	})
	t.Run("PayOrder_get order error_then return error", func(t *testing.T) {
		orderId := "orderId"
		userId := "userId"

		var isActiveOrder *bool
		active := true
		isActiveOrder = &active

		// mock GetOrderById
		amount := 1000
		order := &entity.Order{
			Amount: amount,
		}
		ucTest.transactionRepo.On("GetOrderById", orderId, isActiveOrder).Return(order, errors.New("")).Once()

		// usecase
		err := ucTest.transactionUsecase.PayOrder(&entity.PayOrderRequest{
			OrderId: "orderId",
			Amount:  amount,
			UserId:  userId,
		})

		assert.NotNil(t, err)
	})
	t.Run("PayOrder_insert payment error_then return error", func(t *testing.T) {
		orderId := "orderId"
		userId := "userId"

		var isActiveOrder *bool
		active := true
		isActiveOrder = &active

		// mock GetOrderById
		amount := 1000
		order := &entity.Order{
			Amount: amount,
		}
		ucTest.transactionRepo.On("GetOrderById", orderId, isActiveOrder).Return(order, nil).Once()

		// setup sqlmock for transaction
		db, mockDB, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ucTest.inventoryRepo.On("GetDb").Return(db).Once()
		mockDB.ExpectBegin()

		ucTest.transactionRepo.On("InsertPayment", mock.Anything, mock.Anything).Return(errors.New("")).Once()

		// usecase
		err = ucTest.transactionUsecase.PayOrder(&entity.PayOrderRequest{
			OrderId: "orderId",
			Amount:  amount,
			UserId:  userId,
		})

		assert.NotNil(t, err)
	})
	t.Run("PayOrder_update order error_then return error", func(t *testing.T) {
		orderId := "orderId"
		userId := "userId"

		var isActiveOrder *bool
		active := true
		isActiveOrder = &active

		// mock GetOrderById
		amount := 1000
		order := &entity.Order{
			Amount: amount,
		}
		ucTest.transactionRepo.On("GetOrderById", orderId, isActiveOrder).Return(order, nil).Once()

		// setup sqlmock for transaction
		db, mockDB, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ucTest.inventoryRepo.On("GetDb").Return(db).Once()
		mockDB.ExpectBegin()

		ucTest.transactionRepo.On("InsertPayment", mock.Anything, mock.Anything).Return(nil).Once()
		ucTest.transactionRepo.On("UpdateOrder", mock.Anything, &entity.UpdateOrderRequest{
			OrderId: orderId,
			Status:  entity.OrderStatusSucceeded,
		}).Return(errors.New("")).Once()

		// usecase
		err = ucTest.transactionUsecase.PayOrder(&entity.PayOrderRequest{
			OrderId: "orderId",
			Amount:  amount,
			UserId:  userId,
		})

		assert.NotNil(t, err)
	})
	t.Run("PayOrder_correct payload_then return success", func(t *testing.T) {
		orderId := "orderId"
		productId := "productId"
		userId := "userId"
		warehouseId := "warehouseId"
		quantity := 5

		var isActiveOrder *bool
		active := true
		isActiveOrder = &active

		// mock GetOrderById
		amount := 1000
		order := &entity.Order{
			Amount: amount,
		}
		ucTest.transactionRepo.On("GetOrderById", orderId, isActiveOrder).Return(order, nil).Once()

		// setup sqlmock for transaction
		db, mockDB, err := sqlmock.New()
		assert.NoError(t, err)
		defer db.Close()

		ucTest.inventoryRepo.On("GetDb").Return(db).Once()
		mockDB.ExpectBegin()

		ucTest.transactionRepo.On("InsertPayment", mock.Anything, mock.Anything).Return(nil).Once()
		ucTest.transactionRepo.On("UpdateOrder", mock.Anything, &entity.UpdateOrderRequest{
			OrderId: orderId,
			Status:  entity.OrderStatusSucceeded,
		}).Return(nil).Once()

		// mock GetOrderItemsByOrderId
		orderItems := []*entity.OrderItem{
			{
				ProductId:   productId,
				WarehouseId: warehouseId,
				Quantity:    quantity,
			},
		}
		ucTest.transactionRepo.On("GetOrderItemsByOrderId", orderId).Return(orderItems, nil).Once()

		ucTest.redisRepo.On("InvalidateLockOrderProduct", mock.Anything, mock.Anything).Return(nil).Once()
		ucTest.inventoryRepo.On("InsertProductWarehouse", mock.Anything).Return(nil).Once()

		// usecase
		err = ucTest.transactionUsecase.PayOrder(&entity.PayOrderRequest{
			OrderId: "orderId",
			Amount:  amount,
			UserId:  userId,
		})

		assert.Nil(t, err)
	})
}
