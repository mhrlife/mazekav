package repository

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"math"
	"mazekav/internal/entity"
)

var _ RestaurantRepository = &RestaurantGormRepository{}

type RestaurantGormRepository struct {
	CommonBehaviourRepository[entity.Restaurant]
	db *gorm.DB
}

func NewRestaurantsRepository(db *gorm.DB) *RestaurantGormRepository {
	return &RestaurantGormRepository{
		CommonBehaviourRepository: NewCommonBehaviour[entity.Restaurant](db),
		db:                        db,
	}
}

// Nearby finds restaurants within the radius (in km)
func (r *RestaurantGormRepository) Nearby(ctx context.Context, lat, long, radius float64) ([]entity.Restaurant, error) {
	mbr := `MBRWithin(location, ST_GeomFromText(?, 4326))`
	polygon := generateWKTForMBR(lat, long, radius)
	rad := `ST_Distance_Sphere(ST_GeomFromText(?, 4326),location) < ?`
	point := fmt.Sprintf("Point(%f %f)", lat, long)
	var restaurants []entity.Restaurant
	if err := r.db.WithContext(ctx).Where(mbr, polygon).Where(rad, point, radius*1000).
		Find(&restaurants).Error; err != nil {
		logrus.WithError(err).Errorln("couldn't find nearby restaurants")
		return nil, err
	}
	return restaurants, nil
}

const (
	// Earth's radius in kilometers
	earthRadiusKm = 6371.0
)

// degreesToRadians converts degrees to radians
func degreesToRadians(deg float64) float64 {
	return deg * (math.Pi / 180)
}

// radiansToDegrees converts radians to degrees
func radiansToDegrees(rad float64) float64 {
	return rad * (180 / math.Pi)
}

// calculateBoundingBox calculates the bounding box for a given latitude, longitude, and radius in kilometers
func calculateBoundingBox(lat, lon, radius float64) (minLat, maxLat, minLon, maxLon float64) {
	// Latitude: 1 deg = approx 110.574 kilometers
	// Longitude: 1 deg = approx 111.320*cos(latitude) kilometers
	deltaLat := radius / 110.574
	deltaLon := radius / (111.320 * math.Cos(degreesToRadians(lat)))

	minLat = lat - deltaLat
	maxLat = lat + deltaLat
	minLon = lon - deltaLon
	maxLon = lon + deltaLon

	return minLat, maxLat, minLon, maxLon
}

// generateWKTForMBR generates a WKT polygon for use with ST_GeomFromText in MySQL
func generateWKTForMBR(lat, lon, radius float64) string {
	minLat, maxLat, minLon, maxLon := calculateBoundingBox(lat, lon, radius)

	// Construct Polygon
	return fmt.Sprintf("Polygon((%f %f, %f %f, %f %f, %f %f, %f %f))",
		minLat, minLon,
		minLat, maxLon,
		maxLat, maxLon,
		maxLat, minLon,
		minLat, minLon)
}
