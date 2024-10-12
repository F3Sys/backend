CREATE TYPE node_type AS ENUM ('ENTRY', 'FOODSTALL', 'EXHIBITION');
-- Node table
CREATE TABLE nodes (
    id BIGSERIAL PRIMARY KEY,
    key TEXT UNIQUE,
    otp TEXT UNIQUE,
    name TEXT NOT NULL,
    ip INET UNIQUE,
    type node_type NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Foods table
CREATE TABLE foods (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price INTEGER NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Node Foods table
CREATE TABLE node_foods (
    id BIGSERIAL PRIMARY KEY,
    node_id BIGINT NOT NULL,
    food_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id),
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
    ip INET UNIQUE,
    random INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (model_id) REFERENCES models(id)
);
-- Student table
CREATE TABLE students (
    id BIGSERIAL PRIMARY KEY,
    visitor_id BIGINT UNIQUE NOT NULL,
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
    node_id BIGINT NOT NULL,
    visitor_id BIGINT NOT NULL,
    TYPE entry_logs_type NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id),
    FOREIGN KEY (visitor_id) REFERENCES visitors(id)
);
-- FoodStallLog table
CREATE TABLE food_stall_logs (
    id BIGSERIAL PRIMARY KEY,
    node_food_id BIGINT NOT NULL,
    visitor_id BIGINT NOT NULL,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_food_id) REFERENCES node_foods(id),
    FOREIGN KEY (visitor_id) REFERENCES visitors(id)
);
-- ExhibitionLog table
CREATE TABLE exhibition_logs (
    id BIGSERIAL PRIMARY KEY,
    node_id BIGINT NOT NULL,
    visitor_id BIGINT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id),
    FOREIGN KEY (visitor_id) REFERENCES visitors(id)
);
-- Create indexes
--
-- Visitors Table
CREATE INDEX idx_visitors_ip ON visitors (ip);
CREATE INDEX idx_visitors_id ON visitors (id);
CREATE INDEX idx_visitors_id_random ON visitors (id, random);
-- Nodes Table
CREATE INDEX idx_nodes_ip ON nodes (ip);
CREATE INDEX idx_nodes_key ON nodes (key);
-- EntryLogs Table
CREATE INDEX idx_entry_logs_visitor_id ON entry_logs (visitor_id);
CREATE INDEX idx_entry_logs_node_id ON entry_logs (node_id);
CREATE INDEX idx_entry_logs_type ON entry_logs (type);
CREATE INDEX idx_entry_logs_type_created_at ON entry_logs (type, created_at);
-- FoodStallLogs Table
CREATE INDEX idx_food_stall_logs_node_food_id ON food_stall_logs (node_food_id);
CREATE INDEX idx_food_stall_logs_node_food_id_created_at ON food_stall_logs (node_food_id, created_at);
-- ExhibitionLogs Table
CREATE INDEX idx_exhibition_logs_node_id ON exhibition_logs (node_id);
CREATE INDEX idx_exhibition_logs_node_id_created_at ON exhibition_logs (node_id, created_at);
-- NodeFoods Table
CREATE INDEX idx_node_foods_node_id ON node_foods (node_id);
