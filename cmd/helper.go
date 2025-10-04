package main

import (
	"database/sql"
	"time"
)

func Format(t time.Time) string {
	return t.Format(time.DateOnly)
}

func NFormat(t sql.NullTime) string {
	if !t.Valid {
		return ""
	}

	return t.Time.Format(time.DateOnly)
}

func Today() time.Time {
	return time.Now().Truncate(time.Hour * 24)
}

func SQLInt(n int, b bool) sql.NullInt32 {
	return sql.NullInt32{Int32: int32(n), Valid: b}
}

func SQLTime(t time.Time, b bool) sql.NullTime {
	return sql.NullTime{Time: t, Valid: b}
}
