package repository

import "time"

type UserRepository interface {
	StartTrainig(id int64, startTime time.Time) error
	EndTraining(id int64, endTime time.Time) error
	StartSet(id int64, startTime time.Time) error
	EndSet(id int64, endTime time.Time) error
	SetWeight(id int64, weight float64) error
	SetReps(id int64, reps int) error
	SetExercise(id int64, exercise string) error
	AddExercise(id int64, exercise string) error
	UserCheck(id int64) (bool, error)
	IsExerciseChoosen(id int64) (bool, error)
	RegisterUser(id int64) error
	MaxExerciseId(id int64) (int, error)
	MaxPages(id int64) (int64, error)
	GetPage(id, page int64) ([]string, error)
	IsTrainingActive(id int64) (bool, error)
}

//type AnalitycsRepository interface {
//	ExerciseByWeek(id int64, exercise string)
//}
