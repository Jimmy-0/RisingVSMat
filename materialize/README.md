## Prerequisite
- Docker
- Docker Compose

## How to run
- start all the containers, run `docker compose up -d`
- (optional) check postgres database, run `psql -U postgres -h localhost -p 5432 -d postgres`
- create postgres source for Materialize
    - launch the Materialize CLI, run `docker compose run mzcli`
    - create a Postgres Materialize Source, run
        ``` sql
        CREATE SOURCE IF NOT EXISTS taxi_trips_publication_source FROM POSTGRES
        CONNECTION 'user=postgres port=5432 host=postgres dbname=postgres password=postgres'
        PUBLICATION 'taxi_trips_publication_source';
        ```
    - create a view to represent the upstream publication's original tables, run
        ``` sql
        CREATE MATERIALIZED VIEWS FROM SOURCE taxi_trips_publication_source (taxi_trips);
        ```
    - create a materialized view, run
        ``` sql
        CREATE MATERIALIZED VIEW mv_avg_speed AS
            SELECT COUNT(id) as no_of_trips,
                SUM(distance) as total_distance,
                SUM(duration) as total_duration,
                SUM(distance) / SUM(duration) as avg_speed
            FROM taxi_trips;
        ```
    - query the materialized view, run
        ``` sql
        SELECT * FROM mv_avg_speed;
        ```
