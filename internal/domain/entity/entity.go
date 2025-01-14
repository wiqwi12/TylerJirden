package domain

import "time"

type User struct {
	User_id int64
}

type Set struct {
	user_id  int64
	Exercise string
	Reps     int
	Weight   float64
	Start    time.Time
	End      time.Time
}

type Training struct {
	User_id int64
	Start   time.Time
	End     time.Time
}

type Exercise struct {
	user_id       int64
	exercise_name string
}
