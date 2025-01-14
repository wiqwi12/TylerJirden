package domain

type TrainingsSummary struct {
	TotalWorkOuts          int
	TotalHours             float32
	TotalSets              int
	AverageSetsCount       float32
	AverageExercisesPerSet float32
}

type AverageExercisesStats struct {
	AverageSetsPerTraining float32 //сколько сетов в общем приходилось на одну тренировку
	AverageReps            float32
	AverageWeight          float32
}

type ExerciseStats struct {
	Exercise               string
	AverageReps            float32
	AverageWeight          float32
	AverageSetsPerTraining float32 //сколько сетов определенного упражнения приходилось на тренировку
}
