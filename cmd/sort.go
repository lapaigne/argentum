package main

type Sorting int8

const (
	None Sorting = iota
	Ascending
	Descending
)

type TaskDisplay struct {
	Adrr_obj       Sorting
	Created_date   Sorting
	Until_date     Sorting
	Mark_date      Sorting
	Completed_date Sorting
	Worker         Sorting
}
