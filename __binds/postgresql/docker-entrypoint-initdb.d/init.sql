CREATE TABLE measurements (
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

CREATE TABLE events (
    location_id VARCHAR (255) NOT NULL,
    timestamp INT NOT NULL,
    type VARCHAR (255) NOT NULL,
    PRIMARY KEY (location_id, timestamp, type)
);

CREATE EXTENSION "uuid-ossp";

ALTER TABLE events DROP CONSTRAINT events_pkey;
ALTER TABLE events ADD id uuid DEFAULT uuid_generate_v4 ();
ALTER TABLE events ADD PRIMARY KEY (id);
ALTER TABLE events ADD UNIQUE (location_id, timestamp, type);
ALTER TABLE events RENAME COLUMN timestamp TO start_timestamp;
ALTER TABLE events ADD end_timestamp INT;

ALTER TABLE measurements DROP CONSTRAINT measurements_pkey;
ALTER TABLE measurements ADD id uuid DEFAULT uuid_generate_v4 ();
ALTER TABLE measurements ADD PRIMARY KEY (id);
ALTER TABLE measurements ADD UNIQUE (sensor_id, timestamp);
