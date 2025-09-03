-- This file will run when container starts for the first time
-- Creates the database schema for the log reader application

-- Table to track file uploads
CREATE TABLE IF NOT EXISTS uploads (
    id SERIAL PRIMARY KEY,
    upload_id VARCHAR(255) UNIQUE NOT NULL,
    filename VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    status VARCHAR(50) DEFAULT 'processing'
);

-- Table to store individual log entries from uploaded files
CREATE TABLE IF NOT EXISTS log_entries (
    id SERIAL PRIMARY KEY,
    upload_id VARCHAR(255) REFERENCES uploads(upload_id),
    timestamp VARCHAR(255),
    level VARCHAR(10),
    service VARCHAR(255),
    message TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_uploads_upload_id ON uploads(upload_id);
CREATE INDEX IF NOT EXISTS idx_uploads_status ON uploads(status);
CREATE INDEX IF NOT EXISTS idx_log_entries_upload_id ON log_entries(upload_id);
CREATE INDEX IF NOT EXISTS idx_log_entries_level ON log_entries(level);
CREATE INDEX IF NOT EXISTS idx_log_entries_service ON log_entries(service);
CREATE INDEX IF NOT EXISTS idx_log_entries_created_at ON log_entries(created_at);

-- Insert some sample data for testing (optional)
-- INSERT INTO uploads (upload_id, filename, file_size, status) VALUES 
-- ('sample_001', 'sample.log', 1024, 'completed');
