package integration

import (
	"context"
	"github.com/stretchr/testify/assert"
	"mazekav/internal/entity"
	"mazekav/internal/repository"
	"mazekav/pkg/testutil"
	"testing"
)

func TestRestaurantsNearby(t *testing.T) {
	db := testutil.GetConnection(t)
	assert.NoError(t, db.AutoMigrate(&entity.Restaurant{}))
	ctx := context.Background()

	rr := repository.NewRestaurantsRepository(db)
	assert.NoError(t, rr.Save(ctx, &entity.Restaurant{
		Name:        "CenterGisha",
		Description: "",
		Location:    entity.Location{Latitude: 35.73474551686543, Longitude: 51.377409616420564},
	}))

	assert.NoError(t, rr.Save(ctx, &entity.Restaurant{
		Name:        "CenterAmirAbad",
		Description: "",
		Location:    entity.Location{Latitude: 35.738612180183125, Longitude: 51.39200083310099},
	}))

	assert.NoError(t, rr.Save(ctx, &entity.Restaurant{
		Name:        "CenterShahrAra",
		Description: "",
		Location:    entity.Location{Latitude: 35.72182044179488, Longitude: 51.37088648428088},
	}))

	items, err := rr.Nearby(ctx, 35.73882118359305, 51.37281767472387, 1) // near Gisha
	assert.NoError(t, err)
	assert.Len(t, items, 1)
	assert.Equal(t, "CenterGisha", items[0].Name)

	items, err = rr.Nearby(ctx, 35.728788397778295, 51.37071482290817, 1.3) // near Gisha and ShahrAra
	assert.NoError(t, err)
	assert.Len(t, items, 2)
}
