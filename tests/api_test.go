package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mfawzanid/warehouse-commerce/core/entity"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	email = "email_test@mail.com"

	firstTotalStock      = 100
	updatedTotalStock    = 40
	transferedTotalStock = 35

	pricePerUnit = 1000

	orderedQuantity = 3

	// expecation
	totalAmountToPay = pricePerUnit * orderedQuantity

	sourceTotalStockAfterTransfer = updatedTotalStock - transferedTotalStock
	remainingTotalStockAfterOrder = sourceTotalStockAfterTransfer - orderedQuantity
)

func TestAPI(t *testing.T) {
	if testing.Short() {
		t.Skip("Skip API test")
	}

	testCases := getTestCases()
	ctx := context.Background()
	client := &http.Client{}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			for idx := range tc.Steps {
				step := &tc.Steps[idx]
				request, err := step.Request(t, ctx, &tc)
				request.Header.Set("Content-Type", "application/json")
				request.Header.Set("Accept", "application/json")
				require.NoError(t, err)

				response, err := client.Do(request)

				require.NoError(t, err)
				defer response.Body.Close()

				ReadJSONResult(t, response, step)
				step.Expectation(t, ctx, &tc, response, step.Result)
			}
		})
	}
}

func ReadJSONResult(t *testing.T, resp *http.Response, step *TestCaseStep) {
	var result (map[string]any)
	err := json.NewDecoder(resp.Body).Decode(&result)
	step.Result = result
	require.NoError(t, err)
}

type TestCase struct {
	Name  string
	Steps []TestCaseStep
}

type TestCaseStep struct {
	Request     SetRequestFunc
	Expectation CheckExpectationFunc
	Result      map[string]any
}

type SetRequestFunc func(*testing.T, context.Context, *TestCase) (*http.Request, error)
type CheckExpectationFunc func(*testing.T, context.Context, *TestCase, *http.Response, map[string]any)

const apiURL = "http://localhost:3000"

/*
We test the API using success flow in CreateSuccesTestCaseSteps():
1. Register new user
2. Login using email that registered before
3. Create warehouse, using token in step 2 (Login)
4. Update status warehouse (that created in step 3) become enable
5. Get warehouses, expect only return one warehosue that created in step 3
6. Create shop
7. Get shops
8. Upsert shop to warehouse
9. Create product
10. Update product total stock
11. Create new warehouse as destination warehouse to test transfer product
12. Transfer product to another warehouse
13. Order products
14. Get products by shop to validate total stock after user order product
15. Pay the order
*/
func getTestCases() []TestCase {
	return []TestCase{
		{
			Name:  "Success flow",
			Steps: CreateSuccesTestCaseSteps(),
		},

		// additional test to test the authorization middleware
		{
			Name: "CreateWarehouse_empty token_return unauthorized error",
			Steps: []TestCaseStep{
				{
					Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
						return http.NewRequest("POST", apiURL+"/api/v1/warehouses", nil)
					},
					Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
						require.Equal(t, http.StatusUnauthorized, resp.StatusCode)
					},
				},
			},
		},
	}
}

/*
1. Register new user
2. Login using email that registered before
*/
func RegisterLoginTestCaseStep() []TestCaseStep {
	return []TestCaseStep{
		// 1. Register new user
		{
			Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
				payload := entity.RegisterUserRequest{
					IdentifierType: "email",
					Identifier:     email,
				}

				jsonBody, err := json.Marshal(payload)
				if err != nil {
					return nil, err
				}

				return http.NewRequest("POST", apiURL+"/user/register", bytes.NewReader(jsonBody))
			},
			Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
				require.Equal(t, http.StatusOK, resp.StatusCode)

				tokenStr := data["token"].(string)
				require.NotEmpty(t, tokenStr)
			},
		},
		// 2. Login using email that registered before
		{
			Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
				payload := entity.LoginRequest{
					IdentifierType: "email",
					Identifier:     email,
				}

				jsonBody, err := json.Marshal(payload)
				if err != nil {
					return nil, err
				}

				return http.NewRequest("POST", apiURL+"/user/login", bytes.NewReader(jsonBody))
			},
			Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
				require.Equal(t, http.StatusOK, resp.StatusCode)

				tokenStr := data["token"].(string)
				require.NotEmpty(t, tokenStr)
			},
		},
	}
}

/*
3. Create warehouse, using token in step 2 (Login)
4. Update status warehouse (that created in step 3) become enable
5. Get warehouses, expect only return one warehosue that created in step 3
*/
func CreateWarehouseTestCaseStep() []TestCaseStep {
	return []TestCaseStep{
		// 3. Create warehouse, using token in step 2 (Login)
		{
			Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
				payload := entity.CreateWarehouseRequest{
					Name: "warehouse_test",
				}

				jsonBody, err := json.Marshal(payload)
				if err != nil {
					return nil, err
				}

				// get token from Login response
				loginStep := tc.Steps[1]
				token := loginStep.Result["token"].(string)

				httpReq, _ := http.NewRequest("POST", apiURL+"/api/v1/warehouses", bytes.NewReader(jsonBody))
				httpReq.Header.Set("Authorization", token)
				return httpReq, nil
			},
			Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
				require.Equal(t, http.StatusCreated, resp.StatusCode)
			},
		},
		// 4. Update status warehouse (that created in step 3) become enable
		{
			Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
				// set payload
				payload := map[string]interface{}{
					"enabled": true,
				}
				jsonBody, err := json.Marshal(payload)
				if err != nil {
					return nil, err
				}

				// get warehouse id from CreateWarehouse response
				createWarehouseStep := tc.Steps[2]
				warehouseId := createWarehouseStep.Result["id"].(string)

				// get token from Login response
				loginStep := tc.Steps[1]
				token := loginStep.Result["token"].(string)

				url := fmt.Sprintf("%s/api/v1/warehouses/%s/status", apiURL, warehouseId)
				httpReq, _ := http.NewRequest("PUT", url, bytes.NewReader(jsonBody))
				httpReq.Header.Set("Authorization", token)
				return httpReq, nil
			},
			Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
			},
		},
		// 5. Get warehouses, expect only return one warehosue that created in step
		{
			Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
				// get token from Login response
				loginStep := tc.Steps[1]
				token := loginStep.Result["token"].(string)

				httpReq, _ := http.NewRequest("GET", apiURL+"/api/v1/warehouses?page=1&pageSize=10", nil)
				httpReq.Header.Set("Authorization", token)
				return httpReq, nil
			},
			Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
				require.Equal(t, http.StatusOK, resp.StatusCode)

				// get first warehouse
				warehouses, ok := data["warehouses"].([]interface{})
				require.True(t, ok)
				require.NotEmpty(t, warehouses)

				firstWarehouse, ok := warehouses[0].(map[string]interface{})
				require.True(t, ok)

				// get warehouse id from CreateWarehouse response
				createWarehouseStep := tc.Steps[2]
				warehouseId := createWarehouseStep.Result["id"].(string)

				// validate the warehouse id whether same with warehouseId in step Create Warehouse or not
				require.Equal(t, firstWarehouse["id"].(string), warehouseId)
			},
		},
	}
}

/*
6. Create shop
7. Get shops
8. Upsert shop to warehouse
*/
func CreateShopTestCaseStep() []TestCaseStep {
	return []TestCaseStep{
		// 6. Create shop
		{
			Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
				payload := entity.CreateShopRequest{
					Name: "shop_test",
				}

				jsonBody, err := json.Marshal(payload)
				if err != nil {
					return nil, err
				}

				// get token from Login response
				loginStep := tc.Steps[1]
				token := loginStep.Result["token"].(string)

				httpReq, _ := http.NewRequest("POST", apiURL+"/api/v1/shops", bytes.NewReader(jsonBody))
				httpReq.Header.Set("Authorization", token)
				return httpReq, nil
			},
			Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
				require.Equal(t, http.StatusCreated, resp.StatusCode)

				id, ok := data["id"].(string)
				require.True(t, ok)
				require.NotEmpty(t, id)
			},
		},
		// 7. Get shops
		{
			Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
				// get token from Login response
				loginStep := tc.Steps[1]
				token := loginStep.Result["token"].(string)

				httpReq, _ := http.NewRequest("GET", apiURL+"/api/v1/shops?page=1&pageSize=10", nil)
				httpReq.Header.Set("Authorization", token)
				return httpReq, nil
			},
			Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
				require.Equal(t, http.StatusOK, resp.StatusCode)

				// get first shop
				shops, ok := data["shops"].([]interface{})
				require.True(t, ok)
				require.NotEmpty(t, shops)

				firstShop, ok := shops[0].(map[string]interface{})
				require.True(t, ok)

				// get shop id from CreateShop response
				createShopStep := tc.Steps[5]
				shopId := createShopStep.Result["id"].(string)

				// validate the warehouse id whether same with warehouseId in step Create Warehouse or not
				require.Equal(t, firstShop["id"].(string), shopId)
			},
		},
		// 8. Upsert shop to warehouse
		{
			Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
				// get warehouse id from CreateWarehouse response
				createWarehouseStep := tc.Steps[2]
				warehouseId := createWarehouseStep.Result["id"].(string)

				// get shop id from CreateShop response
				createShopStep := tc.Steps[5]
				shopId := createShopStep.Result["id"].(string)

				// set payload
				req := entity.UpsertShopToWarehousesRequest{
					ShopId:       shopId,
					WarehouseIds: []string{warehouseId},
					Enabled:      true,
				}
				jsonBody, err := json.Marshal(req)
				if err != nil {
					return nil, err
				}

				// get token from Login response
				loginStep := tc.Steps[1]
				token := loginStep.Result["token"].(string)

				url := fmt.Sprintf("%s/api/v1/upsert-shop-warehouses", apiURL)
				httpReq, _ := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
				httpReq.Header.Set("Authorization", token)
				return httpReq, nil
			},
			Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
				require.Equal(t, http.StatusCreated, resp.StatusCode)
			},
		},
	}
}

/*
9. Create product
10. Update product total stock
11. Create new warehouse as destination warehouse to test transfer product
12. Transfer product to another warehouse
*/
func CreateProductTestCaseStep() []TestCaseStep {
	return []TestCaseStep{
		// 9. Create product
		{
			Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
				// get warehouse id from CreateWarehouse response
				createWarehouseStep := tc.Steps[2]
				warehouseId := createWarehouseStep.Result["id"].(string)

				payload := entity.CreateProductRequest{
					Name:        "product_test",
					Price:       pricePerUnit,
					TotalStock:  firstTotalStock,
					WarehouseId: warehouseId,
				}

				jsonBody, err := json.Marshal(payload)
				if err != nil {
					return nil, err
				}

				// get token from Login response
				loginStep := tc.Steps[1]
				token := loginStep.Result["token"].(string)

				httpReq, _ := http.NewRequest("POST", apiURL+"/api/v1/products", bytes.NewReader(jsonBody))
				httpReq.Header.Set("Authorization", token)
				return httpReq, nil
			},
			Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
				require.Equal(t, http.StatusCreated, resp.StatusCode)

				id, ok := data["id"].(string)
				require.True(t, ok)
				require.NotEmpty(t, id)
			},
		},
		// 10. Update product total stock
		{
			Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
				// set payload
				// get warehouse id from CreateWarehouse response
				createWarehouseStep := tc.Steps[2]
				warehouseId := createWarehouseStep.Result["id"].(string)

				payload := map[string]interface{}{
					"totalStock":  updatedTotalStock,
					"warehouseId": warehouseId,
				}
				jsonBody, err := json.Marshal(payload)
				if err != nil {
					return nil, err
				}

				// get product id from CreateProduct response
				createProductStep := tc.Steps[8]
				productId := createProductStep.Result["id"].(string)

				// get token from Login response
				loginStep := tc.Steps[1]
				token := loginStep.Result["token"].(string)

				url := fmt.Sprintf("%s/api/v1/product/%s/stock", apiURL, productId)
				httpReq, _ := http.NewRequest("PUT", url, bytes.NewReader(jsonBody))
				httpReq.Header.Set("Authorization", token)
				return httpReq, nil
			},
			Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
			},
		},
		// 11. Create new warehouse as destination warehouse to test transfer product
		{
			Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
				payload := entity.CreateWarehouseRequest{
					Name: "warehouse_destination_test",
				}

				jsonBody, err := json.Marshal(payload)
				if err != nil {
					return nil, err
				}

				// get token from Login response
				loginStep := tc.Steps[1]
				token := loginStep.Result["token"].(string)

				httpReq, _ := http.NewRequest("POST", apiURL+"/api/v1/warehouses", bytes.NewReader(jsonBody))
				httpReq.Header.Set("Authorization", token)
				return httpReq, nil
			},
			Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
				require.Equal(t, http.StatusCreated, resp.StatusCode)
			},
		},
		// 12. Transfer product to another warehouse
		{
			Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
				// set payload

				// get product id from CreateProduct response
				createProductStep := tc.Steps[8]
				productId := createProductStep.Result["id"].(string)

				// get source warehouse id from CreateWarehouse response in step 3
				createWarehouseStep := tc.Steps[2]
				sourceWarehouseId := createWarehouseStep.Result["id"].(string)

				// get destination warehouse id from CreateWarehouse response in step 11
				createDestinationWarehouseStep := tc.Steps[10]
				destinationWarehouseId := createDestinationWarehouseStep.Result["id"].(string)

				// payload := map[string]interface{}{
				// 	"totalStock":  11,
				// 	"warehouseId": warehouseId,
				// }
				payload := entity.TransferProductRequest{
					ProductId:              productId,
					SourceWarehouseId:      sourceWarehouseId,
					DestinationWarehouseId: destinationWarehouseId,
					TotalStock:             transferedTotalStock,
				}
				jsonBody, err := json.Marshal(payload)
				if err != nil {
					return nil, err
				}

				// get token from Login response
				loginStep := tc.Steps[1]
				token := loginStep.Result["token"].(string)

				url := fmt.Sprintf("%s/api/v1/product/transfer", apiURL)
				httpReq, _ := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
				httpReq.Header.Set("Authorization", token)
				return httpReq, nil
			},
			Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
			},
		},
	}
}

/*
13. Order products
14. Get products by shop to validate total stock after user order product
15. Pay the order
*/
func TransactionTestCaseStep() []TestCaseStep {
	return []TestCaseStep{
		// 13. Order products
		{
			Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
				// set payload

				// get shop id from CreateShop response
				createShopStep := tc.Steps[5]
				shopId := createShopStep.Result["id"].(string)

				// get product id from CreateProduct response
				createProductStep := tc.Steps[8]
				productId := createProductStep.Result["id"].(string)

				payload := map[string]interface{}{
					"items": []interface{}{
						map[string]interface{}{
							"productId": productId,
							"quantity":  orderedQuantity,
						},
					},
				}

				jsonBody, err := json.Marshal(payload)
				if err != nil {
					return nil, err
				}

				// get token from Login response
				loginStep := tc.Steps[1]
				token := loginStep.Result["token"].(string)

				url := fmt.Sprintf("%s/api/v1/shop/%s/order", apiURL, shopId)
				httpReq, _ := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
				httpReq.Header.Set("Authorization", token)
				return httpReq, nil
			},
			Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
				require.Equal(t, http.StatusOK, resp.StatusCode)

				orderId := data["id"].(string)
				require.NotEmpty(t, orderId)
			},
		},
		// 14. Get products by shop to validate total stock after user order product
		{
			Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
				// get shop id from CreateShop response
				createShopStep := tc.Steps[5]
				shopId := createShopStep.Result["id"].(string)

				// get token from Login response
				loginStep := tc.Steps[1]
				token := loginStep.Result["token"].(string)

				url := fmt.Sprintf("%s/api/v1/shops/%s/products?page=1&pageSize=10", apiURL, shopId)
				httpReq, _ := http.NewRequest("GET", url, nil)
				httpReq.Header.Set("Authorization", token)
				return httpReq, nil
			},
			Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
				require.Equal(t, http.StatusOK, resp.StatusCode)

				// get some data to validate the total stock result
				// - product id that create & transfer before
				// - source & destination warehouse id

				// get product id from CreateProduct response
				createProductStep := tc.Steps[8]
				productId := createProductStep.Result["id"].(string)

				// get source warehouse id from CreateWarehouse response in step 3
				createWarehouseStep := tc.Steps[2]
				soureWarehouseId := createWarehouseStep.Result["id"].(string)

				// get products that fetch
				products, ok := data["products"].([]interface{})
				require.True(t, ok)
				require.NotEmpty(t, products)

				for _, productIntf := range products {
					product, ok := productIntf.(map[string]interface{})
					require.True(t, ok)

					// validate total stock expectation for source warehouse
					if product["productId"].(string) == productId {
						if product["warehouseId"] == soureWarehouseId {
							require.Equal(t, remainingTotalStockAfterOrder, int(product["totalStock"].(float64)))
						}
						break
					}
				}
			},
		},
		// 15. Pay order
		{
			Request: func(t *testing.T, ctx context.Context, tc *TestCase) (*http.Request, error) {
				// set payload
				// get order id from Order Product response
				orderProductStep := tc.Steps[12]
				orderId := orderProductStep.Result["id"].(string)
				payload := map[string]interface{}{
					"amount": totalAmountToPay,
				}

				jsonBody, err := json.Marshal(payload)
				if err != nil {
					return nil, err
				}

				// get token from Login response
				loginStep := tc.Steps[1]
				token := loginStep.Result["token"].(string)

				url := fmt.Sprintf("%s/api/v1/order/%s/pay", apiURL, orderId)
				httpReq, _ := http.NewRequest("POST", url, bytes.NewReader(jsonBody))
				httpReq.Header.Set("Authorization", token)
				return httpReq, nil
			},
			Expectation: func(t *testing.T, ctx context.Context, tc *TestCase, resp *http.Response, data map[string]any) {
				require.Equal(t, http.StatusOK, resp.StatusCode)
			},
		},
	}
}

func CreateSuccesTestCaseSteps() []TestCaseStep {
	tcSteps := []TestCaseStep{}
	tcSteps = append(tcSteps, RegisterLoginTestCaseStep()...)
	tcSteps = append(tcSteps, CreateWarehouseTestCaseStep()...)
	tcSteps = append(tcSteps, CreateShopTestCaseStep()...)
	tcSteps = append(tcSteps, CreateProductTestCaseStep()...)
	tcSteps = append(tcSteps, TransactionTestCaseStep()...)

	return tcSteps
}
