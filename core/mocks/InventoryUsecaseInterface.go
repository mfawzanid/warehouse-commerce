// Code generated by mockery v2.50.1. DO NOT EDIT.

package mocks

import (
	entity "mfawzanid/warehouse-commerce/core/entity"

	mock "github.com/stretchr/testify/mock"
)

// InventoryUsecaseInterface is an autogenerated mock type for the InventoryUsecaseInterface type
type InventoryUsecaseInterface struct {
	mock.Mock
}

// CreateProduct provides a mock function with given fields: req
func (_m *InventoryUsecaseInterface) CreateProduct(req *entity.CreateProductRequest) (string, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for CreateProduct")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*entity.CreateProductRequest) (string, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*entity.CreateProductRequest) string); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*entity.CreateProductRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateShop provides a mock function with given fields: req
func (_m *InventoryUsecaseInterface) CreateShop(req *entity.CreateShopRequest) (string, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for CreateShop")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*entity.CreateShopRequest) (string, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*entity.CreateShopRequest) string); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*entity.CreateShopRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// CreateWarehouse provides a mock function with given fields: req
func (_m *InventoryUsecaseInterface) CreateWarehouse(req *entity.CreateWarehouseRequest) (string, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for CreateWarehouse")
	}

	var r0 string
	var r1 error
	if rf, ok := ret.Get(0).(func(*entity.CreateWarehouseRequest) (string, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*entity.CreateWarehouseRequest) string); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Get(0).(string)
	}

	if rf, ok := ret.Get(1).(func(*entity.CreateWarehouseRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetShops provides a mock function with given fields: req
func (_m *InventoryUsecaseInterface) GetShops(req *entity.GetShopsRequest) (*entity.GetShopsResponse, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for GetShops")
	}

	var r0 *entity.GetShopsResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*entity.GetShopsRequest) (*entity.GetShopsResponse, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*entity.GetShopsRequest) *entity.GetShopsResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.GetShopsResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*entity.GetShopsRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetWarehouses provides a mock function with given fields: req
func (_m *InventoryUsecaseInterface) GetWarehouses(req *entity.GetWarehousesRequest) (*entity.GetWarehousesResponse, error) {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for GetWarehouses")
	}

	var r0 *entity.GetWarehousesResponse
	var r1 error
	if rf, ok := ret.Get(0).(func(*entity.GetWarehousesRequest) (*entity.GetWarehousesResponse, error)); ok {
		return rf(req)
	}
	if rf, ok := ret.Get(0).(func(*entity.GetWarehousesRequest) *entity.GetWarehousesResponse); ok {
		r0 = rf(req)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*entity.GetWarehousesResponse)
		}
	}

	if rf, ok := ret.Get(1).(func(*entity.GetWarehousesRequest) error); ok {
		r1 = rf(req)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// TransferProduct provides a mock function with given fields: req
func (_m *InventoryUsecaseInterface) TransferProduct(req *entity.TransferProductRequest) error {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for TransferProduct")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*entity.TransferProductRequest) error); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateProductStock provides a mock function with given fields: req
func (_m *InventoryUsecaseInterface) UpdateProductStock(req *entity.UpdateProductWarehouseTotalStockRequest) error {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for UpdateProductStock")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*entity.UpdateProductWarehouseTotalStockRequest) error); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpdateWarehouseStatus provides a mock function with given fields: req
func (_m *InventoryUsecaseInterface) UpdateWarehouseStatus(req *entity.UpdateWarehouseStatusRequest) error {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for UpdateWarehouseStatus")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*entity.UpdateWarehouseStatusRequest) error); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// UpsertShopToWarehouses provides a mock function with given fields: req
func (_m *InventoryUsecaseInterface) UpsertShopToWarehouses(req *entity.UpsertShopToWarehousesRequest) error {
	ret := _m.Called(req)

	if len(ret) == 0 {
		panic("no return value specified for UpsertShopToWarehouses")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func(*entity.UpsertShopToWarehousesRequest) error); ok {
		r0 = rf(req)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewInventoryUsecaseInterface creates a new instance of InventoryUsecaseInterface. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewInventoryUsecaseInterface(t interface {
	mock.TestingT
	Cleanup(func())
}) *InventoryUsecaseInterface {
	mock := &InventoryUsecaseInterface{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
