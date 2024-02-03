package testutil

import (
	"fmt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"testing"
	"time"
)

// CreateTestDB creates a new in-memory SQLite database for testing purposes.
// It returns a *gorm.DB instance and a cancel function to clean up after the test.
func CreateTestDB(t *testing.T) (*gorm.DB, func()) {
	t.Helper() // Marks the calling function as a test helper function.

	// Use a unique database name for each call to support parallel testing.
	// This ensures each test has its own isolated database environment.
	dbName := fmt.Sprintf("file:memdb%d?mode=memory&cache=shared", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		t.Fatalf("Failed to connect to database: %v", err)
	}

	// Cancel function to clean up database connections
	cancel := func() {
		sqlDB, err := db.DB()
		if err != nil {
			t.Logf("Failed to close database connection: %v", err)
			return
		}
		if err := sqlDB.Close(); err != nil {
			t.Logf("Failed to close database: %v", err)
		}
	}

	return db, cancel
}
