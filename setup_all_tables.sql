-- PostgreSQL Setup Script for GraphQL Validation Tool (Comprehensive)

-- Drop existing database if it exists (for clean setup)
DROP DATABASE IF EXISTS testdb;
DROP USER IF EXISTS testuser;

-- Create the test user
CREATE USER testuser WITH PASSWORD 'testpass';

-- Create the test database
CREATE DATABASE testdb OWNER testuser;

-- Connect to the database
\c testdb

-- =================================================================
-- Basic Tables (users, posts)
-- =================================================================

CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =================================================================
-- Inferred Tables from .gql queries
-- =================================================================

CREATE TABLE clients (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE customers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    credit_info TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE suppliers (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE traders (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE vessels (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE ports (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE specs (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE organizations (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE enquiries (
    id SERIAL PRIMARY KEY,
    status VARCHAR(50),
    port_id INTEGER REFERENCES ports(id),
    order_number VARCHAR(255),
    eta TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE enquiries_customers_onlink (
    id SERIAL PRIMARY KEY,
    enquiry_id INTEGER REFERENCES enquiries(id),
    customer_id INTEGER REFERENCES customers(id)
);

CREATE TABLE enquiry_line_items (
    id SERIAL PRIMARY KEY,
    enquiry_id INTEGER REFERENCES enquiries(id),
    product_id INTEGER REFERENCES products(id),
    quantity INTEGER
);

CREATE TABLE enquiry_line_items_onlink (
    id SERIAL PRIMARY KEY,
    enquiry_id INTEGER REFERENCES enquiries(id),
    product_id INTEGER REFERENCES products(id)
);

CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    enquiry_id INTEGER REFERENCES enquiries(id),
    customer_id INTEGER REFERENCES customers(id),
    seller_id INTEGER REFERENCES users(id),
    status VARCHAR(50),
    date_of_order TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE order_line_items (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id),
    product_id INTEGER REFERENCES products(id),
    price NUMERIC(10, 2)
);

CREATE TABLE order_supplier_links (
    id SERIAL PRIMARY KEY,
    order_id INTEGER REFERENCES orders(id),
    supplier_id INTEGER REFERENCES suppliers(id)
);

CREATE TABLE licenses (
    id SERIAL PRIMARY KEY,
    registration VARCHAR(255),
    supplier_id INTEGER REFERENCES suppliers(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- =================================================================
-- Permissions
-- =================================================================

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO testuser;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO testuser;

-- =================================================================
-- Sample Data
-- =================================================================

INSERT INTO users (name, email) VALUES
    ('John Doe', 'john@example.com'),
    ('Jane Smith', 'jane@example.com'),
    ('Bob Wilson', 'bob@example.com');

INSERT INTO posts (user_id, title, content) VALUES
    (1, 'First Post', 'This is my first post'),
    (1, 'Second Post', 'Another post by John'),
    (2, 'Jane''s Post', 'Hello from Jane'),
    (3, 'Bob''s Thoughts', 'Some thoughts from Bob');

-- =================================================================
-- Verification
-- =================================================================

SELECT 'Setup completed successfully!' as status;
