CREATE TYPE node_type AS ENUM ('ENTRY', 'FOODSTALL', 'EXHIBITION');
-- Foods table
CREATE TABLE foods (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Node table
CREATE TABLE nodes (
    id BIGSERIAL PRIMARY KEY,
    food_id BIGINT,
    key TEXT UNIQUE,
    name TEXT NOT NULL,
    ip INET UNIQUE,
    type node_type NOT NULL,
    is_review BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (food_id) REFERENCES foods(id)
);
-- Battery table
CREATE TABLE batteries (
    id BIGSERIAL PRIMARY KEY,
    node_id BIGINT,
    level INTEGER,
    charging_time INTEGER,
    discharging_time INTEGER,
    charging BOOLEAN,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id)
);
-- Model table
CREATE TABLE models (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Visitor table
CREATE TABLE visitors (
    id BIGSERIAL PRIMARY KEY,
    model_id BIGINT,
    ip INET UNIQUE NOT NULL,
    random INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (model_id) REFERENCES models(id)
);
-- Student table
CREATE TABLE students (
    id BIGSERIAL PRIMARY KEY,
    visitor_id BIGINT UNIQUE,
    grade INTEGER NOT NULL,
    class INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (visitor_id) REFERENCES visitors(id)
);
CREATE TYPE entry_logs_type AS ENUM ('ENTERED', 'LEFT');
-- EntryLog table
CREATE TABLE entry_logs (
    id BIGSERIAL PRIMARY KEY,
    node_id BIGINT,
    visitor_id BIGINT,
    TYPE entry_logs_type NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id),
    FOREIGN KEY (visitor_id) REFERENCES visitors(id)
);
-- FoodStallLog table
CREATE TABLE food_stall_logs (
    id BIGSERIAL PRIMARY KEY,
    node_id BIGINT,
    visitor_id BIGINT,
    food_id BIGINT,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id),
    FOREIGN KEY (visitor_id) REFERENCES visitors(id),
    FOREIGN KEY (food_id) REFERENCES foods(id)
);
-- ExhibitionLog table
CREATE TABLE exhibition_logs (
    id BIGSERIAL PRIMARY KEY,
    node_id BIGINT,
    visitor_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id),
    FOREIGN KEY (visitor_id) REFERENCES visitors(id)
);
-- ExhibitionReviewLog table
CREATE TABLE exhibition_review_logs (
    id BIGSERIAL PRIMARY KEY,
    node_id BIGINT,
    visitor_id BIGINT,
    rating INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id),
    FOREIGN KEY (visitor_id) REFERENCES visitors(id)
);
-- Create indexes
CREATE INDEX idx_nodes_key ON nodes (key);
CREATE INDEX idx_visitors_id ON visitors (id);
CREATE INDEX idx_visitors_id_random ON visitors (id, random);
CREATE INDEX idx_foods_id ON foods (id);