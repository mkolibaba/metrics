package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/mkolibaba/metrics/internal/server/storage"
)

type PostgresStorage struct {
	db *sql.DB
}

func (p *PostgresStorage) GetGauges(ctx context.Context) (map[string]float64, error) {
	stmt, err := p.db.PrepareContext(ctx, "SELECT id, value FROM gauge")
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

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return gauges, nil
}

func (p *PostgresStorage) GetCounters(ctx context.Context) (map[string]int64, error) {
	stmt, err := p.db.PrepareContext(ctx, "SELECT id, delta FROM counter")
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

	if err = rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating rows: %w", err)
	}

	return res, nil
}

func (p *PostgresStorage) GetGauge(ctx context.Context, name string) (float64, error) {
	stmt, err := p.db.PrepareContext(ctx, "SELECT value FROM gauge WHERE id = $1")
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

func (p *PostgresStorage) GetCounter(ctx context.Context, name string) (int64, error) {
	stmt, err := p.db.PrepareContext(ctx, "SELECT delta FROM counter WHERE id = $1")
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

func (p *PostgresStorage) UpdateGauge(ctx context.Context, name string, value float64) (float64, error) {
	err := doUpdateGauge(ctx, p.db, name, value)
	if err != nil {
		return 0, err
	}
	return p.GetGauge(ctx, name)
}

func (p *PostgresStorage) UpdateCounter(ctx context.Context, name string, value int64) (int64, error) {
	err := doUpdateCounter(ctx, p.db, name, value)
	if err != nil {
		return 0, err
	}
	return p.GetCounter(ctx, name)
}

func (p *PostgresStorage) UpdateGauges(ctx context.Context, values []storage.Gauge) error {
	return p.executeInTransaction(ctx, func(tx *sql.Tx) error {
		for _, v := range values {
			err := doUpdateGauge(ctx, tx, v.Name, v.Value)
			if err != nil {
				return fmt.Errorf("update gauge: %w", err)
			}
		}
		return nil
	})
}

func (p *PostgresStorage) UpdateCounters(ctx context.Context, values []storage.Counter) error {
	return p.executeInTransaction(ctx, func(tx *sql.Tx) error {
		for _, v := range values {
			err := doUpdateCounter(ctx, tx, v.Name, v.Value)
			if err != nil {
				return fmt.Errorf("update counter: %w", err)
			}
		}
		return nil
	})
}

func (p *PostgresStorage) executeInTransaction(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := p.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("create tx: %w", err)
	}

	err = fn(tx)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}

	return nil
}

func doUpdateGauge(ctx context.Context, p preparer, name string, value float64) error {
	stmt, err := p.PrepareContext(ctx, `INSERT INTO gauge (id, value) VALUES ($1, $2) 
                              ON CONFLICT (id) DO UPDATE SET value = EXCLUDED.value`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, value)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	return nil
}

func doUpdateCounter(ctx context.Context, p preparer, name string, value int64) error {
	stmt, err := p.PrepareContext(ctx, `INSERT INTO counter (id, delta) VALUES ($1, $2) 
                              ON CONFLICT (id) DO UPDATE SET delta = EXCLUDED.delta + counter.delta`)
	if err != nil {
		return fmt.Errorf("error preparing statement: %w", err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(name, value)
	if err != nil {
		return fmt.Errorf("error executing query: %w", err)
	}

	return nil
}

func New(db *sql.DB) *PostgresStorage {
	return &PostgresStorage{
		db: db,
	}
}

type preparer interface {
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}
