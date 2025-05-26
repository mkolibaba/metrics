package postgres

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/mkolibaba/metrics/internal/server/storage"
)

type PostgresStorage struct {
	db *sql.DB
}

func (p *PostgresStorage) GetGauges() (map[string]float64, error) {
	stmt, err := p.db.Prepare("SELECT id, value FROM gauge")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	gauges := make(map[string]float64)
	for rows.Next() {
		var id string
		var value float64
		err := rows.Scan(&id, &value)
		if err != nil {
			return nil, fmt.Errorf("error reading row: %w", err)
		}
		gauges[id] = value
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return gauges, nil
}

func (p *PostgresStorage) GetCounters() (map[string]int64, error) {
	stmt, err := p.db.Prepare("SELECT id, delta FROM counter")
	if err != nil {
		return nil, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, fmt.Errorf("error executing query: %w", err)
	}

	res := make(map[string]int64)
	for rows.Next() {
		var id string
		var delta int64
		err := rows.Scan(&id, &delta)
		if err != nil {
			return nil, fmt.Errorf("error reading row: %w", err)
		}
		res[id] = delta
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return res, nil
}

func (p *PostgresStorage) GetGauge(name string) (float64, error) {
	stmt, err := p.db.Prepare("SELECT value FROM gauge WHERE id = $1")
	if err != nil {
		return 0, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var value float64
	row := stmt.QueryRow(name)
	err = row.Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrMetricNotFound
		}
		return 0, fmt.Errorf("error reading row: %w", err)
	}

	return value, nil
}

func (p *PostgresStorage) GetCounter(name string) (int64, error) {
	stmt, err := p.db.Prepare("SELECT delta FROM counter WHERE id = $1")
	if err != nil {
		return 0, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	var delta int64
	row := stmt.QueryRow(name)
	err = row.Scan(&delta)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrMetricNotFound
		}
		return 0, fmt.Errorf("error reading row: %w", err)
	}

	return delta, nil
}

func (p *PostgresStorage) UpdateGauge(name string, value float64) (float64, error) {
	stmt, err := p.db.Prepare(`INSERT INTO gauge (id, value) VALUES ($1, $2) 
                              ON CONFLICT (id) DO UPDATE SET value = excluded.value`)
	if err != nil {
		return 0, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, value)
	if err != nil {
		return 0, fmt.Errorf("error executing query: %w", err)
	}

	return p.GetGauge(name)
}

func (p *PostgresStorage) UpdateCounter(name string, value int64) (int64, error) {
	stmt, err := p.db.Prepare(`INSERT INTO counter (id, delta) VALUES ($1, $2) 
                              ON CONFLICT (id) DO UPDATE SET delta = excluded.delta + counter.delta`)
	if err != nil {
		return 0, fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, value)
	if err != nil {
		return 0, fmt.Errorf("error executing query: %w", err)
	}

	return p.GetCounter(name)
}

func New(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}
