package postgres

import (
	"database/sql"
	"errors"
	"github.com/mkolibaba/metrics/internal/server/storage"
	"go.uber.org/zap"
)

type PostgresStorage struct {
	db     *sql.DB
	logger *zap.SugaredLogger
}

func (p *PostgresStorage) GetGauges() map[string]float64 {
	res := make(map[string]float64)
	stmt, err := p.db.Prepare("SELECT id, value FROM gauge")
	if err != nil {
		return res
	}

	rows, err := stmt.Query()
	if err != nil {
		return res
	}

	for rows.Next() {
		var id string
		var value float64
		err := rows.Scan(&id, &value)
		if err != nil {
			return res
		}
		res[id] = value
	}

	if rows.Err() != nil {
		return res
	}

	return res
}

func (p *PostgresStorage) GetCounters() map[string]int64 {
	res := make(map[string]int64)
	stmt, err := p.db.Prepare("SELECT id, delta FROM counter")
	if err != nil {
		return res
	}

	rows, err := stmt.Query()
	if err != nil {
		return res
	}

	for rows.Next() {
		var id string
		var delta int64
		err := rows.Scan(&id, &delta)
		if err != nil {
			return res
		}
		res[id] = delta
	}

	if rows.Err() != nil {
		return res
	}

	return res
}

func (p *PostgresStorage) GetGauge(name string) (float64, error) {
	stmt, err := p.db.Prepare("SELECT value FROM gauge WHERE id = $1")
	if err != nil {
		return 0, err
	}

	var value float64
	row := stmt.QueryRow(name)
	err = row.Scan(&value)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrMetricNotFound
		}
		return 0, err
	}

	return value, nil
}

func (p *PostgresStorage) GetCounter(name string) (int64, error) {
	stmt, err := p.db.Prepare("SELECT delta FROM counter WHERE id = $1")
	if err != nil {
		return 0, err
	}

	var delta int64
	row := stmt.QueryRow(name)
	err = row.Scan(&delta)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, storage.ErrMetricNotFound
		}
		return 0, err
	}

	return delta, nil
}

func (p *PostgresStorage) UpdateGauge(name string, value float64) (float64, error) {
	stmt, err := p.db.Prepare(`INSERT INTO gauge (id, value) VALUES ($1, $2) 
                              ON CONFLICT (id) DO UPDATE SET value = excluded.value`)
	if err != nil {
		return 0, err
	}

	_, err = stmt.Exec(name, value)
	if err != nil {
		return 0, err
	}

	return p.GetGauge(name)
}

func (p *PostgresStorage) UpdateCounter(name string, value int64) (int64, error) {
	stmt, err := p.db.Prepare(`INSERT INTO counter (id, delta) VALUES ($1, $2) 
                              ON CONFLICT (id) DO UPDATE SET delta = excluded.delta + counter.delta`)
	if err != nil {
		return 0, err
	}

	_, err = stmt.Exec(name, value)
	if err != nil {
		return 0, err
	}

	return p.GetCounter(name)
}

func New(db *sql.DB, logger *zap.SugaredLogger) *PostgresStorage {
	return &PostgresStorage{
		db:     db,
		logger: logger,
	}
}
