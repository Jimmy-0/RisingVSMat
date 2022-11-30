# Benchmark

```bash
just build-linux-amd64
```

```bash
for VINSERT in 50 100 500 1000 5000 10000
do
  cd ../risingwave && docker compose up --wait -d && cd -

  RISINGWAVE_HOST="$(docker network inspect risingwave_risingwave | jq  -r '.[].Containers[] | select(.Name == "p4-risingwave") | .IPv4Address | split("/") | .[0]')"

  docker run --rm -it \
    --name risingwave_benchmark \
    --network risingwave_risingwave \
    -v "$(pwd)/bin:/app" \
    ubuntu:20.04 \
    /app/risingwave-benchmark-linux-amd64 \
    --conn-str "host=$RISINGWAVE_HOST port=4566 user=root dbname=dev sslmode=disable" \
    --insert-num $VINSERT \
    --query-factor 0.5 \
    --force-flush \
    --random

  cd ../risingwave && docker compose down && cd -
done
```

```bash
for VINSERT in 50 100 500 1000 5000 10000
do
  cd ../materialize && docker compose up --wait -d && cd -

  MATERIALIZE_HOST="$(docker network inspect materialize_default | jq  -r '.[].Containers[] | select(.Name == "p4-materialized") | .IPv4Address | split("/") | .[0]')"

  docker run --rm -it \
    --name risingwave_benchmark \
    --network materialize_default \
    -v "$(pwd)/bin:/app" \
    ubuntu:20.04 \
    /app/risingwave-benchmark-linux-amd64 \
    --conn-str "host=$MATERIALIZE_HOST port=6875 user=materialize dbname=materialize sslmode=disable" \
    --insert-num $VINSERT \
    --query-factor 0.5 \
    --random

  cd ../materialize && docker compose down && cd -
done
```

Note：
RisingWave has weak support for transactions. Lack of support for some features of postgres, such as CREATE TABLE IF NOT EXIST
