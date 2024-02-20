package integration

import (
	"github.com/stretchr/testify/assert"
	"mazekav/internal/entity"
	"mazekav/pkg/testutil"
	"testing"
)

func TestMigrations(t *testing.T) {
	t.Parallel()
	db := testutil.GetConnection(t)
	err := db.AutoMigrate(
		&entity.User{},
		&entity.Restaurant{},
	)
	assert.NoError(t, err)
	// check indexes
	indexInfo, err := testutil.GetIndexInfo(t, db, "restaurants", "idx_location")
	assert.Equal(t, "SPATIAL", indexInfo.IndexType)
}
