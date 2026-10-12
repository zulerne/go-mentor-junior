//go:build integration

package testdb

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/require"
)

// New applies committed migrations in a schema that is removed after the test.
func New(t *testing.T, migrationsDir string) *pgxpool.Pool {
	t.Helper()
	url := os.Getenv("TEST_DATABASE_URL")
	require.NotEmpty(t, url, "set TEST_DATABASE_URL to a PostgreSQL database allowing CREATE SCHEMA")
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	t.Cleanup(cancel)
	admin, err := pgxpool.New(ctx, url)
	require.NoError(t, err)
	t.Cleanup(admin.Close)
	schema := "store_test_" + strings.ReplaceAll(uuid.NewString(), "-", "")
	_, err = admin.Exec(ctx, "CREATE SCHEMA "+schema)
	require.NoError(t, err)
	t.Cleanup(func() {
		cleanupCtx, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		_, cleanupErr := admin.Exec(cleanupCtx, "DROP SCHEMA "+schema+" CASCADE")
		require.NoError(t, cleanupErr)
	})
	cfg, err := pgxpool.ParseConfig(url)
	require.NoError(t, err)
	cfg.ConnConfig.RuntimeParams["search_path"] = schema
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	require.NoError(t, err)
	t.Cleanup(pool.Close)
	files, err := filepath.Glob(filepath.Join(migrationsDir, "*.up.sql"))
	require.NoError(t, err)
	require.NotEmpty(t, files)
	for _, file := range files {
		sql, readErr := os.ReadFile(file)
		require.NoError(t, readErr)
		_, execErr := pool.Exec(ctx, string(sql))
		require.NoError(t, execErr, "apply migration %s", file)
	}
	return pool
}
