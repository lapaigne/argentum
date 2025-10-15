package main

import (
	"database/sql"
	"time"
)

func PFormat(t time.Time) string {
	return t.Format("02.01.2006")
}

func PNFormat(t sql.NullTime) string {
	if t.Valid {
		return t.Time.Format("02.01.2006")
	}

	return ""
}

func Format(t time.Time) string {
	return t.Format(time.DateOnly)
}

func Status(t sql.NullTime) string {
	if !t.Valid {
		return "✘"
	}

	return "✔"
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

func Diffs(old, upd UserWorker) (dw, du bool) {
	dw = !(old.W.F_name == upd.W.F_name && old.W.I_name == upd.W.I_name && old.W.O_name == upd.W.O_name)
	du = old.U.Level != upd.U.Level
	return dw, du
}
