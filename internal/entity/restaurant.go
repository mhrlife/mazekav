package entity

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Restaurant struct {
	gorm.Model

	Name        string `json:"name"`
	Description string `json:"description"`

	Location Location `json:"location" gorm:"index:idx_location,class:SPATIAL"`
}

func (u Restaurant) Table() string {
	return "restaurants"
}

type Location struct {
	Latitude, Longitude float64
}

func (loc Location) GormDataType() string {
	return "GEOMETRY NOT NULL SRID 4326"
}

func (loc Location) GormValue(ctx context.Context, db *gorm.DB) clause.Expr {
	return clause.Expr{
		SQL: "ST_GeomFromText(?, 4326)",
		Vars: []interface{}{
			fmt.Sprintf("POINT(%f %f)", loc.Latitude, loc.Longitude),
		},
	}
}

// Scan implements the sql.Scanner interface.
func (loc *Location) Scan(src interface{}) error {
	switch src.(type) {
	case []byte:
		var b = src.([]byte)
		if len(b) != 25 {
			return errors.New(fmt.Sprintf("Expected []bytes with length 25, got %d", len(b)))
		}
		var longitude float64
		var latitude float64
		buf := bytes.NewReader(b[9:17])
		err := binary.Read(buf, binary.LittleEndian, &longitude)
		if err != nil {
			return err
		}
		buf = bytes.NewReader(b[17:25])
		err = binary.Read(buf, binary.LittleEndian, &latitude)
		if err != nil {
			return err
		}
		loc.Latitude = latitude
		loc.Longitude = longitude
	default:
		return errors.New(fmt.Sprintf("Expected []byte for Location type, got  %T", src))
	}
	return nil
}
