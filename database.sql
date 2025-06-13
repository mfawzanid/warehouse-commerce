-- This is the SQL script that used to initialize the database schema.

-- users of the commerce platform
CREATE TABLE users (
    id VARCHAR(20) PRIMARY KEY,
    email VARCHAR(50),
    phone_number VARCHAR(50),
    CONSTRAINT unique_email_phone UNIQUE (email, phone_number),
    CONSTRAINT at_least_one_contact CHECK (email IS NOT NULL OR phone_number IS NOT NULL) -- assume just need register email or phone number
);

-- warehouse where products are stocked
CREATE TABLE warehouses (
    id VARCHAR(20) PRIMARY KEY,
    name VARCHAR(100),
    enabled BOOLEAN, 
    UNIQUE(name) -- assumes each warehouse should be unique to prevent confusion
);
CREATE INDEX idx_warehouses_name ON warehouses(name); -- there is need to search by name

-- shops that sell the products
CREATE TABLE shops (
    id VARCHAR(20) PRIMARY KEY,
    name VARCHAR(100),
    UNIQUE(name) -- assume each shop name should be unique to prevent confusion
);
CREATE INDEX idx_shops_name ON shops(name); --there is need to get shop by name

-- a shop can have some warehouses
CREATE TABLE shop_warehouses (
    shop_id VARCHAR(20),
    warehouse_id VARCHAR(20),
    enabled BOOLEAN NOT NULL DEFAULT true,
    UNIQUE(shop_id, warehouse_id),
    CONSTRAINT fk_shop_id FOREIGN KEY (shop_id) REFERENCES shops(id),
    CONSTRAINT fk_warehouse_id FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
);

-- product info
CREATE TABLE products (
    id VARCHAR(20) PRIMARY KEY,
    name VARCHAR(100),
    price INTEGER,
    UNIQUE(name) -- assume each product name should be unique to prevent confusion
);

-- mapping of product stocks per warehouse
CREATE TABLE product_warehouses (
    product_id VARCHAR(20),
    warehouse_id VARCHAR(20),
    total_stock INTEGER NOT NULL DEFAULT 0,
    UNIQUE(product_id, warehouse_id),
    CONSTRAINT fk_pw_product FOREIGN KEY (product_id) REFERENCES products(id),
    CONSTRAINT fk_pw_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
);
CREATE INDEX idx_productid_warehouseid ON product_warehouses(product_id, warehouse_id);

-- user's orders
CREATE TABLE orders (
    id VARCHAR(50) PRIMARY KEY,
    user_id VARCHAR(20) NOT NULL,
    shop_id VARCHAR(20) NOT NULL,
    status VARCHAR(50) NOT NULL, -- pending, succeeded
    amount INTEGER  NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    expired_at TIMESTAMP NOT NULL, -- ssed as fallback to calculate reserved products if redis is unavailable (requires joining order_items)
    CONSTRAINT fk_order_user FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_order_shop FOREIGN KEY (shop_id) REFERENCES shops(id)
);

-- product items per order
CREATE TABLE order_items (
    order_id VARCHAR(50),
    product_id VARCHAR(20),
    shop_id VARCHAR(20),
    warehouse_id VARCHAR(20),
    quantity INTEGER,
    unit_price INTEGER,
    CONSTRAINT fk_item_order FOREIGN KEY (order_id) REFERENCES orders(id),
    CONSTRAINT fk_item_product FOREIGN KEY (product_id) REFERENCES products(id),
    CONSTRAINT fk_item_shop FOREIGN KEY (shop_id) REFERENCES shops(id),
    CONSTRAINT fk_item_warehouse FOREIGN KEY (warehouse_id) REFERENCES warehouses(id)
);
CREATE INDEX idx_order_id ON order_items(order_id);

-- payment for order
CREATE TABLE payments (
    id VARCHAR(50) PRIMARY KEY,
    order_id VARCHAR(50),
    amount INTEGER,
    status VARCHAR(50), -- paid, or another based on next needs
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT fk_payment_order FOREIGN KEY (order_id) REFERENCES orders(id)
);

