CREATE TABLE IF NOT EXITS measurements (
    sensor_id VARCHAR (255) NOT NULL,
    timestamp INT NOT NULL,
    iaq DOUBLE PRECISION,
    co2 DOUBLE PRECISION,
    voc DOUBLE PRECISION,
    pressure DOUBLE PRECISION,
    temperature DOUBLE PRECISION,
    humidity DOUBLE PRECISION,
    PRIMARY KEY (sensor_id, timestamp)
);

CREATE TABLE IF NOT EXITS events (
    location_id VARCHAR (255) NOT NULL,
    timestamp INT NOT NULL,
    type VARCHAR (255) NOT NULL,
    PRIMARY KEY (location_id, timestamp, type)
);