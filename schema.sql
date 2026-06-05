CREATE TABLE IF NOT EXISTS building (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    address TEXT 
);

CREATE TABLE IF NOT EXISTS apartment (
    id SERIAL PRIMARY KEY,
    building_id INTEGER NOT NULL REFERENCES building(id) ON DELETE CASCADE,
    number VARCHAR(50) NOT NULL,
    floor INTEGER NOT NULL,
    sq_meters INTEGER NOT NULL,
    UNIQUE (building_id, number)
);