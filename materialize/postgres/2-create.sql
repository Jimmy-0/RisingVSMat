SET search_path TO taxi;

-- Create tables
DROP TABLE IF EXISTS taxi_trips;

CREATE TABLE taxi_trips
(
    id VARCHAR,
    distance DOUBLE PRECISION,
    duration DOUBLE PRECISION
);

-- Enable REPLICA for tables
ALTER TABLE taxi_trips REPLICA IDENTITY FULL;

-- Create publication on the created tables
CREATE PUBLICATION taxi_trips_publication_source FOR TABLE taxi_trips;