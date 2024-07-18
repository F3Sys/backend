CREATE TYPE node_type AS ENUM ('ENTRY', 'FOODSTALL', 'EXHIBITION');

-- Node table
CREATE TABLE nodes
(
    id         VARCHAR(255) PRIMARY KEY,
    password   VARCHAR(255),
    name       VARCHAR(255) NOT NULL,
    type       node_type NOT NULL,
    price      INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Battery table
CREATE TABLE batteries
(
    id               BIGINT PRIMARY KEY GENERATED always AS IDENTITY,
    node_id          VARCHAR(255),
    level            INTEGER NOT NULL,
    charging_time    INTEGER NOT NULL,
    discharging_time INTEGER NOT NULL,
    charging         BOOLEAN NOT NULL,
    created_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE RESTRICT ON UPDATE
        CASCADE
);

-- Visitor table
CREATE TABLE visitors
(
    id         UUID PRIMARY KEY,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    ip         INET UNIQUE
);

-- Student table
CREATE TABLE students
(
    id         BIGINT PRIMARY KEY GENERATED always AS IDENTITY,
    visitor_id UUID UNIQUE,
    grade      INTEGER NOT NULL,
    class      INTEGER NOT NULL,
    student_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (visitor_id) REFERENCES visitors(id) ON DELETE RESTRICT ON
        UPDATE CASCADE
);

CREATE TYPE entry_logs_type AS ENUM ('ENTERED', 'LEFT');

-- EntryLog table
CREATE TABLE entry_logs
(
    id         BIGINT PRIMARY KEY GENERATED always AS IDENTITY,
    node_id    VARCHAR(255),
    visitor_id UUID,
    TYPE       entry_logs_type NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE RESTRICT ON UPDATE
        CASCADE,
    FOREIGN KEY (visitor_id) REFERENCES visitors(id) ON DELETE RESTRICT ON
        UPDATE CASCADE
);

-- FoodStallLog table
CREATE TABLE food_stall_logs
(
    id         BIGINT PRIMARY KEY GENERATED always AS IDENTITY,
    node_id    VARCHAR(255),
    visitor_id UUID,
    quantity   INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE RESTRICT ON UPDATE
        CASCADE,
    FOREIGN KEY (visitor_id) REFERENCES visitors(id) ON DELETE RESTRICT ON
        UPDATE CASCADE
);

-- ExhibitionLog table
CREATE TABLE exhibition_logs
(
    id         BIGINT PRIMARY KEY GENERATED always AS IDENTITY,
    node_id    VARCHAR(255),
    visitor_id UUID,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (node_id) REFERENCES nodes(id) ON DELETE RESTRICT ON UPDATE
        CASCADE,
    FOREIGN KEY (visitor_id) REFERENCES visitors(id) ON DELETE RESTRICT ON
        UPDATE CASCADE
);