-- This file will run when container starts for the first time

CREATE TABLE uploads (
    id SERIAL PRIMARY KEY,
    upload_id VARCHAR(255) UNIQUE NOT NULL,
    filename VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'processing'
);

CREATE TABLE log_entries (
    id SERIAL PRIMARY KEY,
    upload_id VARCHAR(255) REFERENCES uploads(upload_id),
    timestamp VARCHAR(255),
    level VARCHAR(10),
    service VARCHAR(255),
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);