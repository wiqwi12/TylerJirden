package repository

import (
	domain "GymBot/internal/domain/entity"
	"github.com/xuri/excelize/v2"
	"time"
)

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
	GetMostPopularExercise(id int64) (string, error)
	GetLeastPopularExercise(id int64) (string, error)
	GetAverageWeight(id int64, exercise string) (float64, error)
	GetAverageReps(id int64, exercise string) (int, error)
	GetAverageTrainingsLenght(id int64) (time.Time, error)
	GetTrainingsCount(id int64) (int64, error)
	GetTotalSetsPerExercise(id int64, exercise string) (int64, error)
	GetAverageSetsPerTraining(id int64) (float64, error) //Среднее колличество сэтов за тренировку
	GetTrainings(id int64) ([]domain.Training, error)
	GetExercises(id int64) ([]string, error)
	GetSetsCount(training domain.Training) (int64, error)
	GenerateExelStats(id int64) (excelize.File, error)
	GetAverageExercisesPerTraining(id int64) (float64, error)
}
