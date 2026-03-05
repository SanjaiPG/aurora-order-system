CREATE TABLE inventory (
    product_id SERIAL PRIMARY KEY,
    product_name TEXT,
    stock INT
);

CREATE TABLE users (
    user_id SERIAL PRIMARY KEY,
    username TEXT,
    balance INT
);

CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    user_id INT,
    product_id INT,
    status TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO inventory (product_name, stock)
VALUES ('Laptop', 100);

INSERT INTO users (username, balance)
VALUES ('Jeffrey', 50000);

SELECT * FROM inventory;

SELECT * FROM users;

SELECT * FROM orders;