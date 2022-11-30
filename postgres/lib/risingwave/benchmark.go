package risingwave

import (
	"database/sql"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

type AvgSpeedRow struct {
	trips         int
	totalDistance float64
	totalDuration float64
	avgSpeed      float64
}

type TaxiTripsQuery struct {
	flush                *sql.Stmt
	insertTaxiTripsTable *sql.Stmt
}

func NewTaxiTripsQuery(db *sql.DB) *TaxiTripsQuery {
	query := &TaxiTripsQuery{}
	query.WithInsertTaxiTripsTable(db)
	return query
}

func (q *TaxiTripsQuery) Close() {
	if q.insertTaxiTripsTable != nil {
		q.insertTaxiTripsTable.Close()
	}
	if q.flush != nil {
		q.flush.Close()
	}
}

func MustPrepareStmt(db *sql.DB, query string) *sql.Stmt {
	stmt, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}
	return stmt
}

func (q *TaxiTripsQuery) WithInsertTaxiTripsTable(db *sql.DB) {
	query := `
	INSERT INTO taxi_trips
	VALUES
		($1, $2, $3);
	`
	q.insertTaxiTripsTable = MustPrepareStmt(db, query)
}

func (q *TaxiTripsQuery) WithFlush(db *sql.DB) {
	query := "flush;"
	q.flush = MustPrepareStmt(db, query)
}

type RisingwaveBenchmark struct {
	db          *sql.DB
	verbose     bool
	forceFlush  bool
	insertNum   int
	queryFactor float64
	random      bool
	query       *TaxiTripsQuery
}

func New(connStr string) (*RisingwaveBenchmark, error) {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	return &RisingwaveBenchmark{
		db: db,
	}, nil
}

func (r *RisingwaveBenchmark) WithVerbose(verbose bool) {
	r.verbose = verbose
}

func (r *RisingwaveBenchmark) WithForceFlush(forceFlush bool) {
	r.forceFlush = forceFlush
}

func (r *RisingwaveBenchmark) WithInsertNum(insertNum int) {
	r.insertNum = insertNum
}

func (r *RisingwaveBenchmark) WithQueryFactor(queryFactor float64) {
	r.queryFactor = queryFactor
}

func (r *RisingwaveBenchmark) WithRandom(random bool) {
	r.random = random
}

func (r *RisingwaveBenchmark) CreateTaxiTripsTable() error {
	query := `
	CREATE TABLE taxi_trips(
			id VARCHAR,
			distance DOUBLE PRECISION,
			duration DOUBLE PRECISION
	);
	`
	_, err := r.db.Exec(query)
	return err
}

func (r *RisingwaveBenchmark) CreateAvgSpeedMaterializedView() error {
	query := `
	CREATE MATERIALIZED VIEW mv_avg_speed
	AS
		SELECT COUNT(id) as no_of_trips,
		SUM(distance) as total_distance,
		SUM(duration) as total_duration,
		SUM(distance) / SUM(duration) as avg_speed
		FROM taxi_trips;
	`

	_, err := r.db.Exec(query)
	return err
}

func (r *RisingwaveBenchmark) Benchmark() (err error) {
	err = r.CreateTaxiTripsTable()
	if err != nil {
		return err
	}

	err = r.CreateAvgSpeedMaterializedView()
	if err != nil {
		return err
	}

	r.query = NewTaxiTripsQuery(r.db)
	defer r.query.Close()

	if r.forceFlush {
		r.query.WithFlush(r.db)
	}
	startTime := time.Now()
	err = r.InsertTaxiTripsTable()
	if err != nil {
		return err
	}

	if r.forceFlush {
		_, err = r.query.flush.Exec()
		if err != nil {
			return err
		}
	}

	queryChan := make(chan struct{})
	j := 1
	go func() {
		for i := 0; i < r.insertNum-1; i += 1 {
			err = r.InsertTaxiTripsTable()
			if err != nil {
				log.Fatal(err)
				return
			}
			if r.forceFlush {
				_, err = r.query.flush.Exec()
				if err != nil {
					log.Fatal(err)
					return
				}
			}
			if int(math.Floor((1.0/r.queryFactor)*float64(j))) == i {
				queryChan <- struct{}{}
				j += 1
			}
		}

		close(queryChan)
	}()

	for range queryChan {
		err = r.SelectAvgSpeed()
		if err != nil {
			return err
		}
	}

	endTime := time.Now()
	diffTime := endTime.Sub(startTime)
	fmt.Printf("elapsed time: %d nanoseconds\n", diffTime.Nanoseconds())
	return nil
}

func (r *RisingwaveBenchmark) InsertTaxiTripsTable() error {
	id := uuid.New().String()
	distance := 1.0
	duration := 1.0
	if r.random {
		distance = rand.Float64()
		duration = rand.Float64()
	}
	_, err := r.query.insertTaxiTripsTable.Exec(id, distance, duration)
	return err
}

func (r *RisingwaveBenchmark) SelectAvgSpeed() error {
	avgSpeedRow := &AvgSpeedRow{}
	row := r.db.QueryRow("SELECT no_of_trips, total_distance, total_duration, avg_speed FROM mv_avg_speed;")
	err := row.Scan(&avgSpeedRow.trips, &avgSpeedRow.totalDistance, &avgSpeedRow.totalDuration, &avgSpeedRow.avgSpeed)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	if r.verbose {
		fmt.Printf("trips: %d, totalDistance: %f, totalDuration: %f, avgSpeed: %f\n", avgSpeedRow.trips, avgSpeedRow.totalDistance, avgSpeedRow.totalDuration, avgSpeedRow.avgSpeed)
	}
	return nil
}
