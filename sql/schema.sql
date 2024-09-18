CREATE TYPE node_type AS ENUM ('ENTRY', 'FOODSTALL', 'EXHIBITION');
-- Node table
CREATE TABLE nodes (
    id BIGSERIAL PRIMARY KEY,
    key TEXT UNIQUE,
    name TEXT NOT NULL,
    ip INET UNIQUE,
    type node_type NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Foods table
CREATE TABLE foods (
    id BIGSERIAL PRIMARY KEY,
    node_id BIGINT,
    name TEXT NOT NULL,
    price INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE RESTRICT ON UPDATE CASCADE
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
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
-- Visitor table
CREATE TABLE visitors (
    id BIGSERIAL PRIMARY KEY,
    random INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip INET UNIQUE NOT NULL
);
-- Student table
CREATE TABLE students (
    id BIGSERIAL PRIMARY KEY,
    visitor_id BIGINT UNIQUE,
    grade INTEGER NOT NULL,
    class INTEGER NOT NULL,
    student_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (visitor_id) REFERENCES visitors(id) ON DELETE RESTRICT ON
        UPDATE CASCADE
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
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (visitor_id) REFERENCES visitors(id) ON DELETE RESTRICT ON
        UPDATE CASCADE
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
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (visitor_id) REFERENCES visitors(id) ON DELETE RESTRICT ON
        UPDATE CASCADE,
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
-- ExhibitionLog table
CREATE TABLE exhibition_logs (
    id BIGSERIAL PRIMARY KEY,
    node_id BIGINT,
    visitor_id BIGINT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (visitor_id) REFERENCES visitors(id) ON DELETE RESTRICT ON
        UPDATE CASCADE
);
-- Create indexes
CREATE INDEX idx_nodes_key ON nodes (key);
CREATE INDEX idx_visitors_id ON visitors (id);
CREATE INDEX idx_visitors_id_random ON visitors (id, random);
CREATE INDEX idx_foods_id ON foods (id);