package property

import (
	"time"

	"github.com/lib/pq"
)

type PropertyType string

const (
	PropertyApartment PropertyType = "apartment"
	PropertyVilla     PropertyType = "villa"
	PropertyRoom      PropertyType = "room"
)

type Property struct {
	ID                int64          `db:"id"`
	Name              string         `db:"name"`
	Address           string         `db:"address"`
	Type              PropertyType   `db:"type"`
	BaseRate          float64        `db:"base_rate"`
	MaxGuestsBase     int            `db:"max_guests_base"`
	ExtraRatePerGuest float64        `db:"extra_rate_per_guest"`
	Managers          pq.Int64Array  `db:"managers"`
	Photos            pq.StringArray `db:"photos"`
	Active            bool           `db:"active"`
	CreatedAt         time.Time      `db:"created_at"`
	UpdatedAt         time.Time      `db:"updated_at"`
	CreatedBy         int64          `db:"created_by"`
	UpdatedBy         int64          `db:"updated_by"`
}
