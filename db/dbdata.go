package db

import (
	"time"
)

func DateParser(raw string) *time.Time {
	// NOTE:
	// default YYYY-MM-DD format
	// it is compatible with DATE type in postgres

	time.Parse(time.DateOnly, raw)
	return nil
}
