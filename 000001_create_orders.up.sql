CREATE TABLE orders (
    id SERIAL PRIMARY KEY,
    customer TEXT NOT NULL,
    status TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL
);
