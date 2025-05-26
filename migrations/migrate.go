package migrations

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sort"
)

const AppServer = "server"

func Run(ctx context.Context, db *sql.DB, app string) error {
	switch app {
	case AppServer:
		return doRun(ctx, db, "server")
	}
	return fmt.Errorf("invalid app: %s", app)
}

func doRun(ctx context.Context, db *sql.DB, folder string) error {
	scripts, err := getMigrationScripts(folder)
	if err != nil {
		return err
	}

	sort.Strings(scripts)

	for _, script := range scripts {
		content, err := os.ReadFile(script)
		if err != nil {
			return err
		}

		tx, err := db.BeginTx(ctx, nil)
		if err != nil {
			return err
		}

		_, err = db.ExecContext(ctx, string(content))
		if err != nil {
			return err
		}

		if err := tx.Commit(); err != nil {
			return err
		}
	}

	return nil
}

func getMigrationScripts(scriptsFolder string) ([]string, error) {
	_, f, _, ok := runtime.Caller(0)
	if !ok {
		return []string{}, fmt.Errorf("can not identify path")
	}
	return filepath.Glob(filepath.Join(filepath.Dir(f), scriptsFolder, "*.sql"))
}
