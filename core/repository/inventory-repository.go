package repository

import (
	"database/sql"
	"fmt"
	"mfawzanid/warehouse-commerce/core/entity"
	errorutil "mfawzanid/warehouse-commerce/utils/error"
	"strings"

	"github.com/lib/pq"
)

type InventoryRepositoryInterface interface {
	GetDb() *sql.DB

	// warehouse
	InsertWarehouse(warehouse *entity.Warehouse) error
	UpdateWarehouseStatus(req *entity.UpdateWarehouseStatusRequest) error
	GetWarehouses(req *entity.GetWarehousesRequest) (*entity.GetWarehousesResponse, error)

	// shop
	InsertShop(shop *entity.Shop) error
	GetShops(req *entity.GetShopsRequest) (*entity.GetShopsResponse, error)

	// shop_warehouse
	InsertShopWarehouses(req *entity.UpsertShopToWarehousesRequest) error

	// product
	InsertProduct(tx *sql.Tx, product *entity.Product) error
	InsertProductWarehouse(pw *entity.ProductWarehouse) error
	InsertProductWarehouseTx(tx *sql.Tx, pw *entity.ProductWarehouse) error
	GetProductByName(name string) (*entity.Product, error)
	GetProductDetailsByShopId(req *entity.GetProductDetailsByShopIdRequest) (*entity.GetProductDetailsByShopIdResponse, error)
	GetProductWarehousesByQuery(req *entity.GetProductWarehousesByQueryRequest) ([]*entity.ProductWarehouse, error)
	UpdateProductWarehouseTotalStock(req *entity.UpdateProductWarehouseTotalStockRequest) error
}

type inventoryRepository struct {
	db *sql.DB
}

func NewInventoryRepository(db *sql.DB) InventoryRepositoryInterface {
	return &inventoryRepository{
		db: db,
	}
}

func (r *inventoryRepository) GetDb() *sql.DB {
	return r.db
}

func (r *inventoryRepository) InsertWarehouse(warehouse *entity.Warehouse) error {
	query := `INSERT INTO warehouses (id, name, enabled) VALUES ($1, $2, $3)`

	_, err := r.db.Exec(query, warehouse.Id, warehouse.Name, warehouse.Enabled)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolationErrorCode {
			return errorutil.ErrUniqueViolation
		} else {
			return fmt.Errorf("error repo insert warehouse: %v", err.Error())
		}
	}

	return nil
}

func (r *inventoryRepository) UpdateWarehouseStatus(req *entity.UpdateWarehouseStatusRequest) error {
	if err := req.Validate(); err != nil {
		return err
	}

	query := `UPDATE warehouses SET enabled = $1 WHERE id = $2`

	_, err := r.db.Exec(query, req.Enabled, req.Id)
	if err != nil {
		return fmt.Errorf("error repo update warehouse status: %v", err.Error())
	}

	return nil
}

func (r *inventoryRepository) GetWarehouses(req *entity.GetWarehousesRequest) (*entity.GetWarehousesResponse, error) {
	query := `SELECT id, name, enabled FROM warehouses`

	var conditions []string
	var values []interface{}
	valueIdx := 1

	if req.Enabled != nil {
		conditions = append(conditions, fmt.Sprintf("enabled = $%d", valueIdx))
		values = append(values, *req.Enabled)
		valueIdx++
	}

	if len(req.Ids) > 0 {
		placeholders := make([]string, len(req.Ids))
		for i, id := range req.Ids {
			placeholders[i] = fmt.Sprintf("$%d", valueIdx)
			values = append(values, id)
			valueIdx++
		}
		conditions = append(conditions, fmt.Sprintf("id IN (%s)", strings.Join(placeholders, ",")))
	}

	if req.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name = $%d", valueIdx))
		values = append(values, req.Name)
		valueIdx++
	}

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	if req.Pagination != nil {
		queryCount := fmt.Sprintf(`SELECT COUNT(1) FROM (%s) AS derived`, query)

		err := r.db.QueryRow(queryCount, values...).Scan(&req.Pagination.Total)
		if err != nil {
			return nil, fmt.Errorf("error repo get warehouses: %v", err)
		}

		req.Pagination.SetPagination()

		offset := req.Pagination.GetOffset()
		query = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", query, valueIdx, valueIdx+1)
		values = append(values, req.Pagination.PageSize, offset)
	}

	rows, err := r.db.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("error repo get warehouses: %v", err.Error())
	}
	defer rows.Close()

	var warehouses []*entity.Warehouse
	for rows.Next() {
		warehouse := &entity.Warehouse{}
		err := rows.Scan(&warehouse.Id, &warehouse.Name, &warehouse.Enabled)
		if err != nil {
			return nil, err
		}
		warehouses = append(warehouses, warehouse)
	}

	return &entity.GetWarehousesResponse{
		Warehouses: warehouses,
		Pagination: req.Pagination,
	}, nil
}

func (r *inventoryRepository) InsertShop(shop *entity.Shop) error {
	query := `INSERT INTO shops (id, name) VALUES ($1, $2)`

	_, err := r.db.Exec(query, shop.Id, shop.Name)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolationErrorCode {
			return errorutil.ErrUniqueViolation
		} else {
			return fmt.Errorf("error repo insert shop: %v", err.Error())
		}
	}

	return nil
}

func (r *inventoryRepository) GetShops(req *entity.GetShopsRequest) (*entity.GetShopsResponse, error) {
	query := `SELECT id, name FROM shops`

	var conditions []string
	var values []interface{}
	valueIdx := 1

	if len(req.Ids) > 0 {
		placeholders := make([]string, len(req.Ids))
		for i, id := range req.Ids {
			placeholders[i] = fmt.Sprintf("$%d", valueIdx)
			values = append(values, id)
			valueIdx++
		}
		conditions = append(conditions, fmt.Sprintf("id IN (%s)", strings.Join(placeholders, ",")))
	}
	if req.Name != "" {
		conditions = append(conditions, fmt.Sprintf("name = $%d", valueIdx))
		values = append(values, req.Name)
		valueIdx++
	}

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	if req.Pagination != nil {
		queryCount := fmt.Sprintf(`SELECT COUNT(1) FROM (%s) AS derived`, query)

		err := r.db.QueryRow(queryCount, values...).Scan(&req.Pagination.Total)
		if err != nil {
			return nil, fmt.Errorf("error repo get shops: %v", err)
		}

		req.Pagination.SetPagination()

		offset := req.Pagination.GetOffset()
		query = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", query, valueIdx, valueIdx+1)
		values = append(values, req.Pagination.PageSize, offset)
	}

	rows, err := r.db.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("error repo get shops: %v", err.Error())
	}
	defer rows.Close()

	var shops []*entity.Shop
	for rows.Next() {
		shop := &entity.Shop{}
		err := rows.Scan(&shop.Id, &shop.Name)
		if err != nil {
			return nil, err
		}
		shops = append(shops, shop)
	}

	return &entity.GetShopsResponse{
		Shops:      shops,
		Pagination: req.Pagination,
	}, nil
}

func (r *inventoryRepository) InsertShopWarehouses(req *entity.UpsertShopToWarehousesRequest) error {
	query := `INSERT INTO shop_warehouses (shop_id, warehouse_id, enabled) VALUES %s 
				ON CONFLICT (shop_id, warehouse_id)
				DO UPDATE SET enabled = EXCLUDED.enabled`

	values := []interface{}{}
	placeholders := []string{}

	for i, warehouseId := range req.WarehouseIds {
		placeholders = append(placeholders, fmt.Sprintf("($%d, $%d, $%d)", i*3+1, i*3+2, i*3+3))
		values = append(values, req.ShopId, warehouseId, req.Enabled)
	}

	var queryValues string
	queryValues += strings.Join(placeholders, ", ")

	query = fmt.Sprintf(query, queryValues)

	_, err := r.db.Exec(query, values...)
	if err != nil {
		return fmt.Errorf("error repo insert shop warehouses: %v", err.Error())
	}

	return nil
}

func (r *inventoryRepository) InsertProduct(tx *sql.Tx, product *entity.Product) error {
	query := `INSERT INTO products (id, name, price) VALUES ($1, $2, $3)`

	_, err := tx.Exec(query, product.Id, product.Name, product.Price)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok && pqErr.Code == uniqueViolationErrorCode {
			return errorutil.ErrUniqueViolation
		} else {
			return fmt.Errorf("error repo insert product: %v", err.Error())
		}
	}

	return nil
}

type Execer interface {
	Exec(query string, args ...interface{}) (sql.Result, error)
}

func (r *inventoryRepository) insertProductWarehouse(exec Execer, pw *entity.ProductWarehouse) error {
	if err := pw.Validate(); err != nil {
		return err
	}

	query := `INSERT INTO product_warehouses (product_id, warehouse_id, total_stock) 
				VALUES ($1, $2, $3)
				ON CONFLICT (product_id, warehouse_id)
				DO UPDATE SET total_stock = product_warehouses.total_stock + EXCLUDED.total_stock`

	_, err := exec.Exec(query, pw.ProductId, pw.WarehouseId, pw.TotalStock)
	if err != nil {
		return fmt.Errorf("error repo insert product warehouses: %v", err.Error())
	}

	return nil
}

func (r *inventoryRepository) InsertProductWarehouse(pw *entity.ProductWarehouse) error {
	return r.insertProductWarehouse(r.db, pw)
}

// InsertProductWarehouseTx uses tx for support multiple insert in transfer product between warehouse
func (r *inventoryRepository) InsertProductWarehouseTx(tx *sql.Tx, pw *entity.ProductWarehouse) error {
	return r.insertProductWarehouse(tx, pw)
}

func (r *inventoryRepository) GetProductByName(name string) (*entity.Product, error) {
	query := `SELECT id, name, price FROM products WHERE name = $1`

	product := &entity.Product{}

	err := r.db.QueryRow(query, name).Scan(&product.Id, &product.Name, &product.Price)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorutil.NewErrorCode(errorutil.ErrNotFound, fmt.Errorf("error get product by name '%s' is not found", name))
		} else {
			return nil, fmt.Errorf("error product by name: %v", err.Error())
		}
	}

	return product, nil
}

func (r *inventoryRepository) GetProductDetailsByShopId(req *entity.GetProductDetailsByShopIdRequest) (*entity.GetProductDetailsByShopIdResponse, error) {
	query := `SELECT p.id, p.name, p.price, pw.total_stock, pw.warehouse_id 
				FROM products p
				INNER JOIN product_warehouses pw
				ON p.id = pw.product_id  
				INNER JOIN shop_warehouses sw 
				ON pw.warehouse_id = sw.warehouse_id`

	var conditions []string
	var values []interface{}
	valueIdx := 1

	// mandatory condition
	conditions = append(conditions, "sw.enabled = true")

	conditions = append(conditions, fmt.Sprintf("sw.shop_id = $%d", valueIdx))
	valueIdx++
	values = append(values, req.ShopId)

	if len(req.ProductIds) > 0 {
		placeholders := make([]string, len(req.ProductIds))
		for i, id := range req.ProductIds {
			placeholders[i] = fmt.Sprintf("$%d", valueIdx)
			values = append(values, id)
			valueIdx++
		}
		conditions = append(conditions, fmt.Sprintf("product_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	if req.Pagination != nil {
		queryCount := fmt.Sprintf(`SELECT COUNT(1) FROM (%s) AS derived`, query)

		err := r.db.QueryRow(queryCount, values...).Scan(&req.Pagination.Total)
		if err != nil {
			return nil, fmt.Errorf("error repo get products by shop id: %v", err)
		}

		req.Pagination.SetPagination()

		offset := req.Pagination.GetOffset()
		query = fmt.Sprintf("%s LIMIT $%d OFFSET $%d", query, valueIdx, valueIdx+1)
		values = append(values, req.Pagination.PageSize, offset)
	}

	rows, err := r.db.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("error repo get products by shop id: %v", err.Error())
	}

	var pds []*entity.ProductDetail
	for rows.Next() {
		pd := &entity.ProductDetail{}
		err := rows.Scan(&pd.ProductId, &pd.Name, &pd.Price, &pd.TotalStock, &pd.WarehouseId)
		if err != nil {
			return nil, err
		}
		pds = append(pds, pd)
	}
	defer rows.Close()

	return &entity.GetProductDetailsByShopIdResponse{
		ProductDetails: pds,
		Pagination:     req.Pagination,
	}, nil
}

func (r *inventoryRepository) GetProductWarehousesByQuery(req *entity.GetProductWarehousesByQueryRequest) ([]*entity.ProductWarehouse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}

	query := `SELECT product_id, warehouse_id, total_stock FROM product_warehouses`

	conditions := []string{}
	values := []interface{}{}
	valueIdx := 1

	if len(req.ProductIds) > 0 {
		placeholders := make([]string, len(req.ProductIds))
		for i, id := range req.ProductIds {
			placeholders[i] = fmt.Sprintf("$%d", valueIdx)
			values = append(values, id)
			valueIdx++
		}
		conditions = append(conditions, fmt.Sprintf("product_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(req.WarehouseIds) > 0 {
		placeholders := make([]string, len(req.WarehouseIds))
		for i, id := range req.WarehouseIds {
			placeholders[i] = fmt.Sprintf("$%d", valueIdx)
			values = append(values, id)
			valueIdx++
		}
		conditions = append(conditions, fmt.Sprintf("warehouse_id IN (%s)", strings.Join(placeholders, ",")))
	}

	if len(conditions) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, strings.Join(conditions, " AND "))
	}

	rows, err := r.db.Query(query, values...)
	if err != nil {
		return nil, fmt.Errorf("error repo get product warehouse by query: %v", err.Error())
	}
	defer rows.Close()

	var pws []*entity.ProductWarehouse
	for rows.Next() {
		pw := &entity.ProductWarehouse{}
		err := rows.Scan(&pw.ProductId, &pw.WarehouseId, &pw.TotalStock)
		if err != nil {
			return nil, err
		}
		pws = append(pws, pw)
	}

	return pws, nil
}

func (r *inventoryRepository) UpdateProductWarehouseTotalStock(req *entity.UpdateProductWarehouseTotalStockRequest) (err error) {
	if err := req.Validate(); err != nil {
		return err
	}

	query := `UPDATE product_warehouses SET total_stock = $1 WHERE product_id = $2 AND warehouse_id = $3`

	_, err = r.db.Exec(query, req.TotalStock, req.ProductId, req.WarehouseId)
	if err != nil {
		return fmt.Errorf("error repo update product total stock: %v", err.Error())
	}

	return nil
}
