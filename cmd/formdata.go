package main

import "time"

type FormOpts struct {
	TopOpts []string
	SecOpts []string
	ThrOpts []string
	AddrObj []string
}

type ActorsSelector struct {
}

type FormTempl struct {
	FormOpts FormOpts
	TaskStr  string
	Created  time.Time
	DoUntil  time.Time
	IsUntil  bool
}

type FormData struct {
	TopCat  string
	SecCat  string
	ThrCat  string
	TaskStr string
	AddrObj string
	Created time.Time
	DoUntil time.Time
	IsUntil bool
	Actors  []string
}
