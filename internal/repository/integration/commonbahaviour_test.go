package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"mazekav/internal/repository"
	"mazekav/pkg/testutil"
	"testing"
)

type MyTable struct {
	gorm.Model
	Name string `json:"name"`
}

func (m MyTable) Table() string {
	return "my_tables"
}

func TestCommonBehaviour_Save(t *testing.T) {
	t.Parallel()

	db := testutil.GetConnection(t)

	assert.NoError(t, db.AutoMigrate(&MyTable{}))

	cb := repository.NewCommonBehaviour[MyTable](db)

	myRow := &MyTable{Name: "ali"}
	assert.NoError(t, cb.Save(context.Background(), myRow))
	assert.Equal(t, uint(1), myRow.ID)

	myRow.Name = "Reza"
	assert.NoError(t, cb.Save(context.Background(), myRow))
	assert.Equal(t, uint(1), myRow.ID)
}

func TestCommonBehaviour_ByID(t *testing.T) {
	t.Parallel()
	db := testutil.GetConnection(t)

	assert.NoError(t, db.AutoMigrate(&MyTable{}))

	db.Create(&MyTable{Name: "ali"})
	db.Create(&MyTable{Name: "reza"})

	cb := repository.NewCommonBehaviour[MyTable](db)

	model, err := cb.ByID(context.Background(), 2)
	assert.NoError(t, err)
	assert.Equal(t, "reza", model.Name)
}
