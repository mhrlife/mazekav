package integration

import (
	"github.com/stretchr/testify/assert"
	"mazekav/internal/entity"
	"mazekav/pkg/testutil"
	"testing"
)

func TestRestaurantSpatial(t *testing.T) {
	db := testutil.GetConnection(t)
	assert.NoError(t, db.AutoMigrate(&entity.Restaurant{}))

	restaurant := &entity.Restaurant{
		Name:        "Test Name",
		Description: "Test Desc",
		Location: entity.Location{
			Latitude:  10,
			Longitude: 5,
		},
	}
	assert.NoError(t, db.Save(restaurant).Error)

	var getRestaurant entity.Restaurant
	assert.NoError(t, db.First(&getRestaurant).Error)
	assert.Equal(t, float64(10), getRestaurant.Location.Latitude)
	assert.Equal(t, float64(5), getRestaurant.Location.Longitude)
}
