package models

import (
	"database/sql/driver"
	"time"
)

type SQLiteTimestamp struct {
	Timestamp time.Time
}

// Scan implements the Scanner interface.
func (t *SQLiteTimestamp) Scan(value interface{}) error {
	t.Timestamp = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (t SQLiteTimestamp) Value() (driver.Value, error) {
	return t.Timestamp.Format(time.RFC3339), nil
}

type SQLiteNullTime struct {
	Time  time.Time
	Valid bool
}

// Scan implements the Scanner interface.
func (nt *SQLiteNullTime) Scan(value interface{}) error {
	nt.Time, nt.Valid = value.(time.Time)
	return nil
}

// Value implements the driver Valuer interface.
func (nt SQLiteNullTime) Value() (driver.Value, error) {
	if !nt.Valid {
		return nil, nil
	}
	return nt.Time, nil
}
