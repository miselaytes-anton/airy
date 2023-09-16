CREATE TABLE measurements (
    sensor_id VARCHAR (255) NOT NULL,
    timestamp TIMESTAMP NOT NULL,
    temperature INT,
    humidity INT,
    co2 INT,
    voc INT,
    PRIMARY KEY (sensor_id, timestamp)
);