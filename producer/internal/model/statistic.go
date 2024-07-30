package model

import "time"

type ByHours struct {
	HourStat time.Time `db:"hour" json:"hour_stat"`
	Count    int       `db:"count" json:"count"`
}

type ByDays struct {
	HourStats    map[int]int `json:"hour_stats"`
	OverallCount int
}
