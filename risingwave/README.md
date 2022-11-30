# RisingWave

```bash
docker network create -d bridge risingwave
```

```bash
docker run --rm -it \
  --name risingwave \
  --network risingwave \
  -p 4566:4566 \
  -p 5691:5691 \
  risingwavelabs/risingwave \
  playground
```

```bash
RISINGWAVE_HOST="$(docker inspect risingwave | jq -r ".[].NetworkSettings.Networks.risingwave.IPAddress")"
# OR
docker network inspect risingwave | jq -r '.[].Containers[] | select(.Name == "risingwave") | .IPv4Address'
```

```bash
docker run --rm -it \
  --name psql \
  --network risingwave \
  postgres:15.1 \
  psql -h "$RISINGWAVE_HOST" \
  -p 4566 \
  -d dev \
  -U root
```

```sql
CREATE TABLE taxi_trips(
    id VARCHAR,
    distance DOUBLE PRECISION,
    duration DOUBLE PRECISION
);
```

```sql
CREATE MATERIALIZED VIEW mv_avg_speed
AS
    SELECT COUNT(id) as no_of_trips,
    SUM(distance) as total_distance,
    SUM(duration) as total_duration,
    SUM(distance) / SUM(duration) as avg_speed
    FROM taxi_trips;
```

```sql
INSERT INTO taxi_trips
VALUES
    ('1', 4, 10);
```

```sql
SELECT * FROM mv_avg_speed;
```

```sql
INSERT INTO taxi_trips
VALUES
    ('2', 6, 10);
```

```sql
SELECT * FROM mv_avg_speed;
```