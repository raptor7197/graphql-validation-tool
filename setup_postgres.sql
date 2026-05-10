-- PostgreSQL Setup Script for GraphQL Validation Tool

-- Drop existing database if it exists (for clean setup)
DROP DATABASE IF EXISTS testdb;
DROP USER IF EXISTS testuser;

-- Create the test user
CREATE USER testuser WITH PASSWORD 'testpass';

-- Create the test database
CREATE DATABASE testdb OWNER testuser;

-- Connect to the database
\c testdb

-- Create the users table
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create the posts table
CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id),
    title VARCHAR(255) NOT NULL,
    content TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Grant permissions to testuser
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO testuser;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO testuser;

-- Insert sample data into users
INSERT INTO users (name, email) VALUES
    ('John Doe', 'john@example.com'),
    ('Jane Smith', 'jane@example.com'),
    ('Bob Wilson', 'bob@example.com');

-- Insert sample data into posts
INSERT INTO posts (user_id, title, content) VALUES
    (1, 'First Post', 'This is my first post'),
    (1, 'Second Post', 'Another post by John'),
    (2, 'Jane''s Post', 'Hello from Jane'),
    (3, 'Bob''s Thoughts', 'Some thoughts from Bob');

-- Verify the data
SELECT 'Users table:' as info;
SELECT * FROM users;

SELECT 'Posts table:' as info;
SELECT * FROM posts;

SELECT 'Setup completed successfully!' as status;
