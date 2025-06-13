# Warehouse Commerce API

## Overview
This API server implementing a simplified e-commerce that stored the products in some warehouses.

## Functionality
- Do simple register & login using phone number or email
- Create a warehouse, get the list, and update the status
- Create a shop 
- Bind a shop to some warehouses
- Create a product and set it in some warehouses
- Update product stock in a warehouse
- Transfer product from one warehouse to another
- Get products in a shop
- Order products with stock reservation
- Do payment with stock updating

## APIs
Some APIs that we need to cover all functionality requirements:
### User Domain
- Register User
- Login

### Inventory Domain
- Create Warehouse
- Get Warehouses
- Update Warehouse Status
- Create Shop
- Get Shops
- Upsert (Bind) Shop To Warehouses
- Create Product
- Update Product Stock
- Transfer Product

### Transaction Domain
- Get Products in a Shop
- Order Products
- Pay Order


## Database Schema
This schema supports warehouse-commerce platform with users, shops, warehouses, products, orders, and payments.

---

### **users**
Stores user information including contact details.

| Column        | Type        | Constraints                                                          | Description                                |
|---------------|-------------|----------------------------------------------------------------------|--------------------------------------------|
| id            | VARCHAR(20) | PRIMARY KEY                                                          | Unique user ID                              |
| email         | VARCHAR(50) |                                                                      | User email (optional)                       |
| phone_number  | VARCHAR(50) |                                                                      | User phone number (optional)                |

**Constraints**:
- At least one of `email` or `phone_number` must be provided.
- Combination of `email` and `phone_number` must be unique.

---

### **warehouses**
Stores warehouse metadata.

| Column  | Type        | Constraints        | Description                         |
|---------|-------------|--------------------|-------------------------------------|
| id      | VARCHAR(20) | PRIMARY KEY        | Unique warehouse ID                 |
| name    | VARCHAR(100)| UNIQUE             | Warehouse name                      |
| enabled | BOOLEAN     |                    | Whether the warehouse is active     |

---

### **shops**
Stores shop metadata.

| Column  | Type        | Constraints        | Description                     |
|---------|-------------|--------------------|---------------------------------|
| id      | VARCHAR(20) | PRIMARY KEY        | Unique shop ID                  |
| name    | VARCHAR(100)| UNIQUE             | Shop name                       |

---

### **shop_warehouses**
Links shops with their warehouses.

| Column        | Type        | Constraints                                            | Description                                |
|---------------|-------------|--------------------------------------------------------|--------------------------------------------|
| shop_id       | VARCHAR(20) | FOREIGN KEY → shops(id)                                | Shop ID                                     |
| warehouse_id  | VARCHAR(20) | FOREIGN KEY → warehouses(id)                           | Warehouse ID                                |
| enabled       | BOOLEAN     | NOT NULL DEFAULT true                                  | Whether this shop-warehouse binding is active |

**Unique constraint**: `(shop_id, warehouse_id)`

---

### **products**
Stores product details.

| Column | Type        | Constraints        | Description                 |
|--------|-------------|--------------------|-----------------------------|
| id     | VARCHAR(20) | PRIMARY KEY        | Unique product ID           |
| name   | VARCHAR(100)| UNIQUE             | Product name                |
| price  | INTEGER     |                    | Product price               |

---

### **product_warehouses**
Tracks stock levels of products in warehouses.

| Column       | Type        | Constraints                                            | Description                          |
|--------------|-------------|--------------------------------------------------------|--------------------------------------|
| product_id   | VARCHAR(20) | FOREIGN KEY → products(id)                             | Product ID                           |
| warehouse_id | VARCHAR(20) | FOREIGN KEY → warehouses(id)                           | Warehouse ID                         |
| total_stock  | INTEGER     | NOT NULL DEFAULT 0                                     | Stock available in this warehouse    |

**Unique constraint**: `(product_id, warehouse_id)`  

---

### **orders**
Stores customer orders.

| Column     | Type        | Constraints                            | Description                             |
|------------|-------------|----------------------------------------|-----------------------------------------|
| id         | VARCHAR(50) | PRIMARY KEY                            | Unique order ID                          |
| user_id    | VARCHAR(20) | FOREIGN KEY → users(id)                | User that order                            |
| shop_id    | VARCHAR(20) | FOREIGN KEY → shops(id)                | Shop that received the order             |
| status     | VARCHAR(50) | NOT NULL                               | Order status (e.g., pending, succeeded)  |
| amount     | INTEGER     | NOT NULL                               | Total price amount                      |
| created_at | TIMESTAMP   | NOT NULL DEFAULT CURRENT_TIMESTAMP     | Order creation time                      |
| expired_at | TIMESTAMP   | NOT NULL                               | Used as fallback to expire reservations |

---

### **order_items**
Details of products in an order.

| Column       | Type        | Constraints                                          | Description                                |
|--------------|-------------|------------------------------------------------------|--------------------------------------------|
| order_id     | VARCHAR(50) | FOREIGN KEY → orders(id)                             | Order ID                                   |
| product_id   | VARCHAR(20) | FOREIGN KEY → products(id)                           | Product ID                                 |
| shop_id      | VARCHAR(20) | FOREIGN KEY → shops(id)                              | Shop ID                                    |
| warehouse_id | VARCHAR(20) | FOREIGN KEY → warehouses(id)                         | Source warehouse                           |
| quantity     | INTEGER     |                                                      | Quantity ordered                           |
| unit_price   | INTEGER     |                                                      | Price per item at the time of order        |

---

### **payments**
Stores payment info related to orders.

| Column     | Type        | Constraints                            | Description                              |
|------------|-------------|----------------------------------------|------------------------------------------|
| id         | VARCHAR(50) | PRIMARY KEY                            | Payment ID                                |
| order_id   | VARCHAR(50) | FOREIGN KEY → orders(id)               | Associated order                          |
| amount     | INTEGER     |                                        | Amount paid                               |
| status     | VARCHAR(50) |                                        | Payment status (e.g., paid)      |
| created_at | TIMESTAMP   | NOT NULL DEFAULT CURRENT_TIMESTAMP     | Time the payment was recorded             |

---

### **Relationships**
- A `user` places an `order` from a `shop`
- A `shop` operates through one or more `warehouses`
- A `product` is stocked in one or more `warehouses`
- An `order` contains multiple `order_items`
- `payments` are linked to `orders`


## Initiate The Project
Run the following command to initialize the project:
```
make init
```

## Running the API Server
Use Docker to run the project:
```
docker compose up --build
```

The API server will be accessible at http://localhost:3000

If you modify the database schema in `database.sql`, you must reinitialize the database by running:
```
docker compose down --volumes
```

## Testing

### Unit test
Run the unit tests:
```
make test
```

### API Test
Run the API server:
```
docker compose up --build
```
Then run API tests:
```
make test_api
```
