package repository

import (
	"context"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
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

	db, cancel := testutil.CreateTestDB(t)
	defer cancel()

	assert.NoError(t, db.AutoMigrate(&MyTable{}))

	cb := NewCommonBehaviour[MyTable](db)

	myRow := &MyTable{Name: "ali"}
	assert.NoError(t, cb.Save(context.Background(), myRow))
	assert.Equal(t, uint(1), myRow.ID)

	myRow.Name = "Reza"
	assert.NoError(t, cb.Save(context.Background(), myRow))
	assert.Equal(t, uint(1), myRow.ID)
}

func TestCommonBehaviour_ByID(t *testing.T) {
	db, cancel := testutil.CreateTestDB(t)
	defer cancel()

	assert.NoError(t, db.AutoMigrate(&MyTable{}))

	db.Create(&MyTable{Name: "ali"})
	db.Create(&MyTable{Name: "reza"})

	cb := NewCommonBehaviour[MyTable](db)

	model, err := cb.ByID(context.Background(), 2)
	assert.NoError(t, err)
	assert.Equal(t, "reza", model.Name)
}
