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
	)
	assert.NoError(t, err)
}
