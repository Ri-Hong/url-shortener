-- Create the URLs table
CREATE TABLE urls (
    id SERIAL PRIMARY KEY,
    short_code VARCHAR(255) UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    click_count INT DEFAULT 0,
    created_at TIMESTAMP DEFAULT NOW()
);
