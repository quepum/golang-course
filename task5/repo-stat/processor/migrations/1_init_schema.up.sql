CREATE TABLE IF NOT EXISTS repositories (
    id SERIAL PRIMARY KEY,
    owner VARCHAR(255) NOT NULL,
    repo VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    description TEXT,
    stars INT,
    forks INT,
    created_at TIMESTAMP WITH TIME ZONE,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
    UNIQUE(owner, repo)
    );