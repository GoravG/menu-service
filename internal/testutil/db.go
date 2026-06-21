package testutil

import (
	"database/sql"
	"os"
	"strings"
	"testing"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"restaurant-menu-api/internal/db"
)

func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dsn := os.Getenv("TEST_DB_DSN")
	if dsn == "" {
		t.Skip("TEST_DB_DSN not set, skipping integration test")
	}

	ensureTestDatabase(t, dsn)

	database, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("failed to open test database: %v", err)
	}

	database.SetMaxOpenConns(5)
	database.SetMaxIdleConns(2)
	database.SetConnMaxLifetime(5 * time.Minute)

	if err := database.Ping(); err != nil {
		t.Fatalf("failed to ping test database: %v", err)
	}

	db.CreateTablesIfNotExists(database)

	t.Cleanup(func() {
		truncateTables(database)
		database.Close()
	})

	return database
}

func ensureTestDatabase(t *testing.T, dsn string) {
	t.Helper()

	dbName := dbNameFromDSN(dsn)
	if dbName == "" {
		t.Fatal("TEST_DB_DSN must include a database name")
	}

	adminDSN := adminDSNFromDSN(dsn)
	adminDB, err := sql.Open("mysql", adminDSN)
	if err != nil {
		t.Fatalf("failed to open admin database connection: %v", err)
	}
	defer adminDB.Close()

	if err := adminDB.Ping(); err != nil {
		t.Fatalf("failed to ping admin database connection: %v", err)
	}

	_, err = adminDB.Exec("CREATE DATABASE IF NOT EXISTS `" + dbName + "`")
	if err != nil {
		t.Fatalf("failed to create test database %q: %v", dbName, err)
	}
}

func truncateTables(database *sql.DB) {
	tables := []string{
		"menu_tags_list",
		"menu_price_lists",
		"menu_items",
		"tags",
		"categories",
	}

	database.Exec("SET FOREIGN_KEY_CHECKS = 0")
	for _, table := range tables {
		database.Exec("TRUNCATE TABLE " + table)
	}
	database.Exec("SET FOREIGN_KEY_CHECKS = 1")
}

func dbNameFromDSN(dsn string) string {
	withoutParams := strings.SplitN(dsn, "?", 2)[0]
	idx := strings.LastIndex(withoutParams, "/")
	if idx == -1 || idx == len(withoutParams)-1 {
		return ""
	}
	return withoutParams[idx+1:]
}

func adminDSNFromDSN(dsn string) string {
	dbName := dbNameFromDSN(dsn)
	withoutParams := strings.SplitN(dsn, "?", 2)[0]
	adminBase := strings.TrimSuffix(withoutParams, "/"+dbName) + "/"
	if params := strings.SplitN(dsn, "?", 2); len(params) == 2 {
		return adminBase + "?" + params[1]
	}
	return adminBase
}

func LinkMenuTag(t *testing.T, database *sql.DB, menuItemName, tag string) {
	t.Helper()

	_, err := database.Exec(
		"INSERT INTO menu_tags_list (menu_item_name, tag) VALUES (?, ?)",
		menuItemName, tag,
	)
	if err != nil {
		t.Fatalf("failed to link tag %q to menu item %q: %v", tag, menuItemName, err)
	}
}
