-- Create node_type as a text column with constraints
CREATE TABLE nodes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    key TEXT UNIQUE,
    name TEXT NOT NULL,
    ip TEXT UNIQUE,
    type TEXT CHECK(type IN ('ENTRY', 'FOODSTALL', 'EXHIBITION')) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- Foods table
CREATE TABLE foods (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    node_id INTEGER,
    name TEXT NOT NULL,
    price INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
-- Battery table
CREATE TABLE batteries (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    node_id INTEGER,
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
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    random INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip TEXT UNIQUE NOT NULL
);
-- Student table
CREATE TABLE students (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    visitor_id INTEGER UNIQUE,
    grade INTEGER NOT NULL,
    class INTEGER NOT NULL,
    student_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (visitor_id) REFERENCES visitors(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
-- Create entry_logs_type as a text column with constraints
CREATE TABLE entry_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    node_id INTEGER,
    visitor_id INTEGER,
    type TEXT CHECK(type IN ('ENTERED', 'LEFT')) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (visitor_id) REFERENCES visitors(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
-- FoodStallLog table
CREATE TABLE food_stall_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    node_id INTEGER,
    visitor_id INTEGER,
    food_id INTEGER,
    quantity INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (visitor_id) REFERENCES visitors(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (food_id) REFERENCES foods(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
-- ExhibitionLog table
CREATE TABLE exhibition_logs (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    node_id INTEGER,
    visitor_id INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (visitor_id) REFERENCES visitors(id) ON DELETE RESTRICT ON UPDATE CASCADE
);
-- Create indexes
CREATE INDEX idx_nodes_key ON nodes (key);
CREATE INDEX idx_visitors_id ON visitors (id);
CREATE INDEX idx_visitors_id_random ON visitors (id, random);
CREATE INDEX idx_foods_id ON foods (id);
-- Enable foreign key constraints
PRAGMA foreign_keys = ON;