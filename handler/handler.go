package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"

	"mfawzanid/warehouse-commerce/core/entity"
	"mfawzanid/warehouse-commerce/core/usecase"
	"mfawzanid/warehouse-commerce/generated"
	errorutil "mfawzanid/warehouse-commerce/utils/error"
	generalutil "mfawzanid/warehouse-commerce/utils/general"
)

type handler struct {
	userUsecase        usecase.UserUsecaseInterface
	inventoryUsecase   usecase.InventoryUsecaseInterface
	transactionUsecase usecase.TransactionUsecaseInterface
}

func NewServer(userUsecase usecase.UserUsecaseInterface, inventoryUsecase usecase.InventoryUsecaseInterface, transactionUsecase usecase.TransactionUsecaseInterface) *handler {
	return &handler{
		userUsecase:        userUsecase,
		inventoryUsecase:   inventoryUsecase,
		transactionUsecase: transactionUsecase,
	}
}

func (h *handler) GetHealth(ctx echo.Context) error {
	return ctx.JSON(http.StatusOK, generalutil.MapAny{
		"status": "ok",
	})
}

func (h *handler) RegisterUser(ctx echo.Context) error {
	var req entity.RegisterUserRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
			errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
		})
	}

	token, err := h.userUsecase.RegisterUser(&req)
	if err != nil {
		switch errorutil.GetErrorType(err) {
		case errorutil.ErrBadRequest:
			return ctx.JSON(http.StatusBadRequest, generalutil.MapAny{
				errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusBadRequest, errorutil.GetOriginalError(err)),
			})
		default:
			if ctx != nil {
				return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
					errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
				})
			}
			log.Printf("error register user: %v\n", errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err))
			return nil
		}
	}

	return ctx.JSON(http.StatusOK, generated.RegisterUserResponse{
		Token: token,
	})
}

func (h *handler) Login(ctx echo.Context) error {
	var req entity.LoginRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
			errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
		})
	}

	token, err := h.userUsecase.Login(&req)
	if err != nil {
		switch errorutil.GetErrorType(err) {
		case errorutil.ErrBadRequest:
			return ctx.JSON(http.StatusBadRequest, generalutil.MapAny{
				errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusBadRequest, errorutil.GetOriginalError(err)),
			})
		case errorutil.ErrNotFound:
			return ctx.JSON(http.StatusNotFound, generalutil.MapAny{
				errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusNotFound, errorutil.GetOriginalError(err)),
			})
		default:
			if ctx != nil {
				return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
					errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
				})
			}
			log.Printf("error login: %v\n", errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err))
			return nil
		}
	}

	return ctx.JSON(http.StatusOK, generated.LoginResponse{
		Token: token,
	})
}

func (h *handler) CreateWarehouse(ctx echo.Context) error {
	var req entity.CreateWarehouseRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
			errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
		})
	}

	warehouseId, err := h.inventoryUsecase.CreateWarehouse(&req)
	if err != nil {
		log.Printf("error create warehouse: %v\n", errorutil.GetOriginalError(err))

		switch errorutil.GetErrorType(err) {
		case errorutil.ErrBadRequest:
			return ctx.JSON(http.StatusBadRequest, generalutil.MapAny{
				errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusBadRequest, errorutil.GetOriginalError(err)),
			})
		default:
			if ctx != nil {
				return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
					errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
				})
			}
			log.Printf("error create warehouse: %v\n", errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err))
			return nil
		}
	}

	return ctx.JSON(http.StatusCreated, generated.CreateWarehouseResponse{
		Id: warehouseId,
	})
}

func (h *handler) UpdateWarehouseStatus(ctx echo.Context, warehouseId string) error {
	var req entity.UpdateWarehouseStatusRequest
	req.Id = warehouseId

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
			errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
		})
	}

	if err := h.inventoryUsecase.UpdateWarehouseStatus(&req); err != nil {
		switch errorutil.GetErrorType(err) {
		case errorutil.ErrBadRequest:
			return ctx.JSON(http.StatusBadRequest, generalutil.MapAny{
				errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusBadRequest, errorutil.GetOriginalError(err)),
			})
		default:
			if ctx != nil {
				return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
					errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
				})
			}
			log.Printf("error update warehouse status: %v\n", errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err))
			return nil
		}
	}

	return ctx.JSON(http.StatusOK, generalutil.MapAny{
		errorutil.Message: "Successfully updated the warehouse status",
	})
}

func (h *handler) GetWarehouses(ctx echo.Context, req generated.GetWarehousesParams) error {
	pagination := entity.ParseToPagination(req.Page, req.PageSize)

	resp, err := h.inventoryUsecase.GetWarehouses(&entity.GetWarehousesRequest{
		Pagination: pagination,
	})
	if err != nil {
		switch errorutil.GetErrorType(err) {
		default:
			if ctx != nil {
				return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
					errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
				})
			}
			log.Printf("error get warehouses: %v\n", errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err))
			return nil
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h *handler) CreateShop(ctx echo.Context) error {
	var req entity.CreateShopRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
			errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
		})
	}

	shopId, err := h.inventoryUsecase.CreateShop(&req)
	if err != nil {
		switch errorutil.GetErrorType(err) {
		case errorutil.ErrBadRequest:
			return ctx.JSON(http.StatusBadRequest, generalutil.MapAny{
				errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusBadRequest, errorutil.GetOriginalError(err)),
			})
		default:
			if ctx != nil {
				return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
					errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
				})
			}
			log.Printf("error create shop: %v\n", errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err))
			return nil
		}
	}

	return ctx.JSON(http.StatusCreated, generated.CreateShopResponse{
		Id: shopId,
	})
}

func (h *handler) GetShops(ctx echo.Context, req generated.GetShopsParams) error {
	pagination := entity.ParseToPagination(req.Page, req.PageSize)

	resp, err := h.inventoryUsecase.GetShops(&entity.GetShopsRequest{
		Pagination: pagination,
	})
	if err != nil {
		switch errorutil.GetErrorType(err) {
		default:
			if ctx != nil {
				return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
					errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
				})
			}
			log.Printf("error get shops: %v\n", errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err))
			return nil
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h *handler) UpsertShopToWarehouses(ctx echo.Context) error {
	var req entity.UpsertShopToWarehousesRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
			errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
		})
	}

	err := h.inventoryUsecase.UpsertShopToWarehouses(&req)
	if err != nil {
		switch errorutil.GetErrorType(err) {
		case errorutil.ErrBadRequest:
			return ctx.JSON(http.StatusBadRequest, generalutil.MapAny{
				errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusBadRequest, errorutil.GetOriginalError(err)),
			})
		default:
			if ctx != nil {
				return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
					errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
				})
			}
			log.Printf("error upsert shop to warehouses: %v\n", errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err))
			return nil
		}
	}

	return ctx.JSON(http.StatusCreated, generalutil.MapAny{
		errorutil.Message: "Successfully insert or update shop to warehouse binding",
	})
}

func (h *handler) CreateProduct(ctx echo.Context) error {
	var req entity.CreateProductRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
			errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
		})
	}

	productId, err := h.inventoryUsecase.CreateProduct(&req)
	if err != nil {
		switch errorutil.GetErrorType(err) {
		case errorutil.ErrBadRequest:
			return ctx.JSON(http.StatusBadRequest, generalutil.MapAny{
				errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusBadRequest, errorutil.GetOriginalError(err)),
			})
		default:
			if ctx != nil {
				return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
					errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
				})
			}
			log.Printf("error create product: %v\n", errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err))
			return nil
		}
	}

	return ctx.JSON(http.StatusCreated, generated.CreateProductResponse{
		Id: productId,
	})
}

func (h *handler) UpdateProductStock(ctx echo.Context, productId string) error {
	var req entity.UpdateProductWarehouseTotalStockRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
			errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
		})
	}

	req.ProductId = productId

	err := h.inventoryUsecase.UpdateProductStock(&req)
	if err != nil {
		switch errorutil.GetErrorType(err) {
		case errorutil.ErrBadRequest:
			return ctx.JSON(http.StatusBadRequest, generalutil.MapAny{
				errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusBadRequest, errorutil.GetOriginalError(err)),
			})
		default:
			if ctx != nil {
				return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
					errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
				})
			}
			log.Printf("error update product stock: %v\n", errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err))
			return nil
		}
	}

	return ctx.JSON(http.StatusOK, generalutil.MapAny{
		errorutil.Message: "Product's stock is updated",
	})
}

func (h *handler) TransferProduct(ctx echo.Context) error {
	var req entity.TransferProductRequest

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
			errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
		})
	}

	err := h.inventoryUsecase.TransferProduct(&req)
	if err != nil {
		switch errorutil.GetErrorType(err) {
		case errorutil.ErrBadRequest:
			return ctx.JSON(http.StatusBadRequest, generalutil.MapAny{
				errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusBadRequest, errorutil.GetOriginalError(err)),
			})
		default:
			if ctx != nil {
				return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
					errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
				})
			}
			log.Printf("error transfer product: %v\n", errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err))
			return nil
		}
	}

	return ctx.JSON(http.StatusOK, generalutil.MapAny{
		errorutil.Message: "Transfer product is completed",
	})
}

func (h *handler) GetProductsByShopId(ctx echo.Context, shopId string, params generated.GetProductsByShopIdParams) error {
	pagination := entity.ParseToPagination(params.Page, params.PageSize)

	resp, err := h.transactionUsecase.GetProductDetailsByShopId(&entity.GetProductDetailsByShopIdRequest{
		ShopId:     shopId,
		Pagination: pagination,
	})
	if err != nil {
		switch errorutil.GetErrorType(err) {
		case errorutil.ErrBadRequest:
			return ctx.JSON(http.StatusBadRequest, generalutil.MapAny{
				errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusBadRequest, errorutil.GetOriginalError(err)),
			})
		default:
			if ctx != nil {
				return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
					errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
				})
			}
			log.Printf("error get products by shop id: %v\n", errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err))
			return nil
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}

func (h *handler) OrderProducts(ctx echo.Context, shopId string) error {
	userIdInterface := ctx.Get(entity.ContextUserId)
	userId, ok := userIdInterface.(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, generalutil.MapAny{
			errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, fmt.Errorf("order products: invalid user id in context")),
		})
	}

	var req entity.OrderProductsRequest
	req.UserId = userId
	req.ShopId = shopId

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
			errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
		})
	}

	orderId, err := h.transactionUsecase.OrderProducts(&req)
	if err != nil {
		switch errorutil.GetErrorType(err) {
		case errorutil.ErrBadRequest:
			return ctx.JSON(http.StatusBadRequest, generalutil.MapAny{
				errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusBadRequest, errorutil.GetOriginalError(err)),
			})
		default:
			if ctx != nil {
				return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
					errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
				})
			}
			log.Printf("error order products: %v\n", errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err))
			return nil
		}
	}

	return ctx.JSON(http.StatusOK, generated.OrderProductsResponse{
		Id: orderId,
	})
}

func (h *handler) PayOrder(ctx echo.Context, orderId string) error {
	userIdInterface := ctx.Get(entity.ContextUserId)
	userId, ok := userIdInterface.(string)
	if !ok {
		return ctx.JSON(http.StatusUnauthorized, generalutil.MapAny{
			errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, fmt.Errorf("order products: invalid user id in context")),
		})
	}

	req := entity.PayOrderRequest{
		UserId:  userId,
		OrderId: orderId,
	}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
			errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
		})
	}

	err := h.transactionUsecase.PayOrder(&req)
	if err != nil {
		switch errorutil.GetErrorType(err) {
		case errorutil.ErrBadRequest:
			return ctx.JSON(http.StatusBadRequest, generalutil.MapAny{
				errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusBadRequest, errorutil.GetOriginalError(err)),
			})
		default:
			if ctx != nil {
				return ctx.JSON(http.StatusInternalServerError, generalutil.MapAny{
					errorutil.Error: errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err),
				})
			}
			log.Printf("error pay order: %v\n", errorutil.CombineHTTPErrorMessage(http.StatusInternalServerError, err))
			return nil
		}
	}

	return ctx.JSON(http.StatusOK, generalutil.MapAny{
		errorutil.Message: "Payment is completed",
	})
}
