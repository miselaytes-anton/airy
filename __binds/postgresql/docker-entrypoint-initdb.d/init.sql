CREATE TABLE measurements (
    sensor_id VARCHAR (255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    iaq DOUBLE PRECISION,
    co2 DOUBLE PRECISION,
    voc DOUBLE PRECISION,
    pressure DOUBLE PRECISION,
    temperature DOUBLE PRECISION,
    humidity DOUBLE PRECISION,
    PRIMARY KEY (sensor_id, timestamp)
);