package postgres

import (
	"GymBot/internal"
	domain "GymBot/internal/domain/entity"
	"GymBot/internal/domain/repository"
	"database/sql"
	"fmt"
	"log"
	"log/slog"
	"sort"
	"time"

	"github.com/xuri/excelize/v2"

	"github.com/Masterminds/squirrel"
)

type UserRepositoryDB struct {
	Db *sql.DB
}

func NewUserRepositoryDb(db *sql.DB) repository.UserRepository {
	return &UserRepositoryDB{
		Db: db,
	}
}

func (u *UserRepositoryDB) StartTrainig(id int64, startTime time.Time) error {

	q := squirrel.Insert("trainings").Columns("user_id", "start_time").Values(id, startTime).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()

	if err != nil {
		slog.Error("Start Training ToSql error:", err)
		return err
	}

	_, err = u.Db.Query(query, args...)
	if err != nil {
		slog.Error("Start Trainig Query Error:", err)
		return err
	}

	return nil
}

func (u *UserRepositoryDB) EndTraining(id int64, endTime time.Time) error {

	q := squirrel.Update("trainings").Set("end_time", endTime).Where(
		squirrel.And{
			squirrel.Eq{"user_id": id},
			squirrel.Expr("end_time IS NULL"),
		}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("End Training ToSql error:", err)
		return err
	}

	_, err = u.Db.Query(query, args...)
	if err != nil {
		slog.Error("End training error:", err)
		return err
	}

	return nil
}

func (u *UserRepositoryDB) StartSet(id int64, startTime time.Time) error {

	q := squirrel.Update("sets").Set("start_time", startTime).Where(
		squirrel.And{
			squirrel.Eq{"user_id": id},
			squirrel.Expr("start_time IS NULL AND exercise_name IS NOT NULL"),
		}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("Start Set ToSql Error:", err)
		return err
	}

	_, err = u.Db.Query(query, args...)
	if err != nil {
		slog.Error("Start Set Query Error:", err)
		return err
	}

	return nil
}

func (u *UserRepositoryDB) EndSet(id int64, endTime time.Time) error {

	q := squirrel.Update("sets").Set("end_time", endTime).Where(
		squirrel.And{
			squirrel.Eq{"user_id": id},
			squirrel.Expr("end_time IS NULL"),
		}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("end Set ToSql Error:", err)
		return err
	}

	_, err = u.Db.Query(query, args...)
	if err != nil {
		slog.Error("end Set Query Error:", err)
		return err
	}

	return nil
}

func (u *UserRepositoryDB) SetWeight(id int64, weight float64) error {

	q := squirrel.Update("sets").Set("weight", weight).Where(
		squirrel.And{
			squirrel.Eq{"user_id": id},
			squirrel.Expr("weight IS NULL"),
		}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("set weight ToSql Error:", err)
		return err
	}

	_, err = u.Db.Query(query, args...)
	if err != nil {
		slog.Error("set weight Query Error:", err)
		return err
	}

	return nil
}

func (u *UserRepositoryDB) SetReps(id int64, reps int) error {

	q := squirrel.Update("sets").Set("reps", reps).Where(
		squirrel.And{
			squirrel.Eq{"user_id": id},
			squirrel.Expr("reps IS NULL"),
		}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("set reps ToSql Error:", err)
		return err
	}

	_, err = u.Db.Query(query, args...)
	if err != nil {
		slog.Error("set reps Query Error:", err)
		return err
	}

	return nil
}

func (u *UserRepositoryDB) AddExercise(id int64, exercise string) error {

	q := squirrel.Insert("exercises").Columns("name", "user_id").Values(exercise, id).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("add exercise ToSql Error:", err)
		return err
	}

	_, err = u.Db.Exec(query, args...)
	if err != nil {
		slog.Error("add exercise Query Error:", err)
		return err
	}

	return nil
}

func (u *UserRepositoryDB) SetExercise(id int64, exercise string) error {

	q := squirrel.Insert("sets").Columns("exercise_name", "user_id").Values(exercise, id).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("Set Exercise ToSql Error:", err)
		return err
	}

	_, err = u.Db.Query(query, args...)
	if err != nil {
		slog.Error("Set Exercise Query Error:", err)
		return err
	}

	return nil
}

func (u *UserRepositoryDB) IsExerciseChoosen(id int64) (bool, error) {

	q := squirrel.Select("COUNT(*)").From("sets").Where(
		squirrel.And{
			squirrel.Eq{"user_id": id},
			squirrel.Expr("end_time IS NULL AND exercise_name IS NOT NULL"),
		}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("Is exercise choosen ToSql error:", err)
		return false, err
	}

	var count int
	err = u.Db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		slog.Error("Is exercise choosen QueryRow error:", err)
		return false, err
	}

	return count > 0, nil
}

func (u *UserRepositoryDB) RegisterUser(id int64) error {

	q := squirrel.Insert("users").Columns("user_id").Values(id).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("Register user ToSql Error:", err)
		return err
	}

	_, err = u.Db.Exec(query, args...)
	if err != nil {
		slog.Error("Register user Exec Error:", err)
		return err
	}

	return nil

}

func (u *UserRepositoryDB) UserCheck(id int64) (bool, error) {

	q := squirrel.
		Select("COUNT(*)").
		From("users").
		Where(squirrel.Eq{"user_id": id}).
		PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("User check ToSql error:", err)
		return false, err
	}

	var count int
	err = u.Db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		slog.Error("User check QueryRow error:", err)
		return false, err
	}

	return count > 0, nil
}

func (u *UserRepositoryDB) MaxExerciseId(id int64) (int, error) {

	increment := squirrel.Select("MAX(exercise_id)").From("exercises").Where(
		squirrel.Eq{"user_id": id}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := increment.ToSql()
	if err != nil {
		slog.Error("increment add exercise ToSql Error:", err)
		return 0, err
	}

	var max int

	err = u.Db.QueryRow(query, args...).Scan(&max)
	if err != nil {
		slog.Error("increment add exercise QueryRow Error:", err)
		return 0, err
	}

	return max, nil
}

func (u *UserRepositoryDB) GetPage(id, page int64) ([]string, error) {

	maxId := page * 5
	minId := maxId - 5

	q := squirrel.Select("name").
		From("exercises").
		Where(
			squirrel.And{
				squirrel.Eq{"user_id": id},
				squirrel.LtOrEq{"exercise_id": maxId},
				squirrel.Gt{"exercise_id": minId},
			},
		).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("GetPage ToSql Error:", err)
		return nil, err
	}

	rows, err := u.Db.Query(query, args...)
	if err != nil {
		slog.Error("GetPage Query Error:", err)
		return nil, err
	}

	defer rows.Close()

	var exercises []string
	for rows.Next() {
		var exercise string
		if err := rows.Scan(&exercise); err != nil {
			log.Fatal(err)
		}
		exercises = append(exercises, exercise)
	}

	return exercises, nil
}

func (u *UserRepositoryDB) MaxPages(id int64) (int64, error) {

	var count int64

	q := squirrel.Select("count(*)").From("exercises").Where(squirrel.Eq{"user_id": id}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("MaxPages ToSql Error:", err)
		return 0, err
	}

	slog.Info("SQL Query", slog.String("Query", query), slog.Any("Args", args))

	row := u.Db.QueryRow(query, args...)
	if err := row.Scan(&count); err != nil {
		slog.Error("MaxPages QueryRow Error:", err)
		return 0, err
	}

	return (count + 4) / 5, nil
}

func (u *UserRepositoryDB) IsTrainingActive(id int64) (bool, error) {

	q := squirrel.Select("COUNT(*)").
		From("trainings").
		Where(
			squirrel.And{
				squirrel.Eq{"user_id": id},
				squirrel.Expr("start_time IS NOT NULL"),
				squirrel.Expr("end_time IS NULL"),
			}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("IsTrainingActive ToSql error:", err)
		return false, err
	}

	var count int

	err = u.Db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		slog.Error("IsTrainingActive QueryRow error:", err)
		return false, err
	}

	return count > 0, nil
}

func (u *UserRepositoryDB) GetMostPopularExercise(id int64) (string, error) {

	q := squirrel.Select("exercise_name").
		From("sets").
		Where(squirrel.Eq{"user_id": id}).
		GroupBy("exercise_name").
		OrderBy("COUNT(*) DESC").
		Limit(1).
		PlaceholderFormat(squirrel.Dollar)

	var exercise string

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("GetMostPopularExercise ToSql Error:", err)
		return "", err
	}

	row := u.Db.QueryRow(query, args...)

	err = row.Scan(&exercise)
	if err != nil {
		slog.Error("GetMostPopularExercise QueryRow Error:", err)
		return "", err
	}

	return exercise, nil

}
func (u *UserRepositoryDB) GetLeastPopularExercise(id int64) (string, error) {

	q := squirrel.Select("exercise_name").
		From("sets").
		Where(squirrel.Eq{"user_id": id}).
		GroupBy("exercise_name").
		OrderBy("COUNT(*) ASC").
		Limit(1).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("GetMostPopularExercise ToSql Error:", err)
		return "", err
	}

	var exercise string

	err = u.Db.QueryRow(query, args...).Scan(&exercise)
	if err != nil {
		slog.Error("GetMostPopularExercise QueryRow Error:", err)
		return "", err
	}

	return exercise, nil
}

func (u *UserRepositoryDB) GetAverageWeight(id int64, exercise string) (float64, error) {
	q := squirrel.Select("AVG(weight)").From("sets").Where(
		squirrel.And{
			squirrel.Eq{"user_id": id},
			squirrel.Eq{"exercise_name": exercise},
		}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("GetAverageWeight ToSql Error:", err)
		return 0, err
	}

	var weight sql.NullFloat64

	err = u.Db.QueryRow(query, args...).Scan(&weight)
	if err != nil {
		slog.Error("GetAverageWeight QueryRow Error:", err)
		return 0, err
	}

	var result float64

	if weight.Valid {
		result = weight.Float64
	} else {
		result = 0.0
	}

	return result, nil
}
func (u *UserRepositoryDB) GetAverageReps(id int64, exercise string) (string, error) {

	q := squirrel.Select("AVG(reps)").From("sets").Where(squirrel.Eq{
		"user_id":       id,
		"exercise_name": exercise,
	}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("GetAverageReps ToSql Error:", err)
		return "", err
	}

	var reps sql.NullString

	err = u.Db.QueryRow(query, args...).Scan(&reps)
	if err != nil {
		slog.Error("GetAverageReps QueryRow Error:", err)
		return "", err
	}

	var result string

	if reps.Valid {
		result = reps.String[0:3]
	} else {
		result = "0"
	}

	return result, nil

}

func (u *UserRepositoryDB) GetAverageTrainingsLenght(id int64) (time.Duration, error) {

	trainings, err := u.GetTrainings(id)
	if err != nil {
		slog.Warn("GetAverageTrainingsLenght GetTrainings Error:", err)
		return time.Duration(0), err
	}

	var avgLength time.Duration
	var durationsSlice []time.Duration

	for _, training := range trainings {
		duration := training.End.Sub(training.Start)
		durationsSlice = append(durationsSlice, duration)
	}

	var total time.Duration

	for _, v := range durationsSlice {
		total += v
	}

	if len(durationsSlice) > 0 {
		avgLength = total / time.Duration(len(durationsSlice)) // Делим на time.Duration(len)
	}

	return avgLength, nil

}

func (u *UserRepositoryDB) GetTrainingsCount(id int64) (int64, error) {
	q := squirrel.Select("COUNT(*)").From("trainings").Where(
		squirrel.Eq{"user_id": id}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("GetTrainingsCount ToSql Error:", err)
		return 0, err
	}

	var count int64

	err = u.Db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		slog.Error("GetTrainingsCount QueryRow Error:", err)
		return 0, err
	}

	return count, nil

}

func (u *UserRepositoryDB) GetTotalSetsPerExercise(id int64, exercise string) (int64, error) {

	q := squirrel.Select("COUNT(*)").From("sets").Where(
		squirrel.And{
			squirrel.Eq{"user_id": id},
			squirrel.Eq{"exercise_name": exercise},
		}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("GetTotalSetsPerExercise ToSql Error:", err)
		return 0, err
	}

	var totalSets int64

	err = u.Db.QueryRow(query, args...).Scan(&totalSets)
	if err != nil {
		slog.Error("GetTotalSetsPerExercise QueryRow Error:", err)
		return 0, err
	}

	return totalSets, nil
}

func (u *UserRepositoryDB) GetTrainings(id int64) ([]domain.Training, error) {
	var trainings []domain.Training

	q := squirrel.Select("user_id", "start_time", "end_time").From("trainings").Where(
		squirrel.Eq{"user_id": id}).OrderBy("start_time").PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("GetTrainingsByUserID ToSql Error:", err)
		return nil, err
	}

	rows, err := u.Db.Query(query, args...)
	if err != nil {
		slog.Error("GetTrainingsByUserID Query Error:", err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var training domain.Training
		if err := rows.Scan(
			&training.User_id,
			&training.Start,
			&training.End,
		); err != nil {
			slog.Error("GetTrainingsByUserID Scan Error:", err)
			return nil, err
		}
		trainings = append(trainings, training)
	}

	if err := rows.Err(); err != nil {
		slog.Error("GetTrainingsByUserID Rows Error:", err)
		return nil, err
	}

	return trainings, nil
}

func (u *UserRepositoryDB) GetExercises(id int64) ([]string, error) {

	q := squirrel.Select("name").From("exercises").Where(
		squirrel.Eq{"user_id": id}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("GetExercises ToSql Error:", err)
		return nil, err
	}

	rows, err := u.Db.Query(query, args...)
	if err != nil {
		slog.Error("GetExercises Query Error:", err)
		return nil, err
	}
	defer rows.Close()

	var exercises []string
	for rows.Next() {
		var exerciseName string
		if err := rows.Scan(&exerciseName); err != nil {
			slog.Error("GetExercises Scan Error:", err)
			return nil, err
		}
		exercises = append(exercises, exerciseName)
	}

	if err := rows.Err(); err != nil {
		slog.Error("GetExercises Rows Error:", err)
		return nil, err
	}

	return exercises, nil

}

func (u *UserRepositoryDB) GetSetsCount(training domain.Training, exercise string) (int, error) {

	q := squirrel.Select("COUNT(*)").From("sets").Where(squirrel.And{
		squirrel.Eq{"user_id": training.User_id},
		squirrel.GtOrEq{"start_time": training.Start},
		squirrel.LtOrEq{"end_time": training.End},
		squirrel.Eq{"exercise_name": exercise},
	}).PlaceholderFormat(squirrel.Dollar)

	query, args, err := q.ToSql()
	if err != nil {
		slog.Error("GetSetsCount ToSql Error:", err)
		return 0, err
	}

	var count int

	err = u.Db.QueryRow(query, args...).Scan(&count)
	if err != nil {
		slog.Error("GetSetsCount QueryRow Error:", err)
		return 0, err
	}

	return count, nil

}

func (u *UserRepositoryDB) GetAverageSetsPerExerise(id int64, exercise string) (string, error) {

	trainings, err := u.GetTrainings(id)
	if err != nil {
		slog.Error("GetTrainings Error:", err)
		return "", err
	}

	var count []int

	for _, training := range trainings {
		val, err := u.GetSetsCount(training, exercise)
		if err != nil {
			slog.Error("GetSetsCount Error:", err)
			return "", err
		}
		count = append(count, val)
	}

	var result int

	for _, v := range count {
		result += v
	}

	r := fmt.Sprintf("%.2f", float64(result)/float64(len(count)))

	return r, nil

}

func (u *UserRepositoryDB) GetAverageExercisesPerTraining(id int64) (float64, error) {

	trainings, err := u.GetTrainings(id)
	if err != nil {
		slog.Error("GetTrainings Error:", err)
		return 0, err
	}

	var count []int64

	for _, training := range trainings {

		q := squirrel.Select("COUNT(DISTINCT exercise_name)").From("sets").Where(
			squirrel.And{
				squirrel.Eq{"user_id": training.User_id},
				squirrel.GtOrEq{"start_time": training.Start},
				squirrel.LtOrEq{"end_time": training.End},
			}).PlaceholderFormat(squirrel.Dollar)

		query, args, err := q.ToSql()
		if err != nil {
			slog.Error("Avg exercises per training err: ", err)
			return 0, err
		}

		var uniqueExercises int64

		err = u.Db.QueryRow(query, args...).Scan(&uniqueExercises)
		if err != nil {
			slog.Error("GetExercises QueryRow Error:", err)
			return 0, err
		}

		count = append(count, uniqueExercises)
	}

	var result float64
	for _, v := range count {
		result += float64(v)
	}

	return result / float64(len(count)), nil

}

func (u *UserRepositoryDB) GetAverageSetsPerTraining(id int64) (float64, error) {

	trainings, err := u.GetTrainings(id)
	if err != nil {
		slog.Error("GetTrainings Error:", err)
		return 0, err
	}

	var count []int64

	for _, training := range trainings {
		q := squirrel.Select("COUNT(*)").From("sets").Where(
			squirrel.And{
				squirrel.Eq{"user_id": training.User_id},
				squirrel.GtOrEq{"start_time": training.Start},
				squirrel.LtOrEq{"end_time": training.End},
			}).PlaceholderFormat(squirrel.Dollar)

		query, args, err := q.ToSql()
		if err != nil {
			slog.Error("GetSetsCount Error:", err)
			return 0, err
		}

		var c int64
		err = u.Db.QueryRow(query, args...).Scan(&c)
		if err != nil {
			slog.Error("GetSetsCount QueryRow Error:", err)
			return 0, err
		}

		count = append(count, c)

	}

	var result float64

	for _, v := range count {
		result += float64(v)
	}

	return result / float64(len(count)), nil
}

func (u *UserRepositoryDB) GenerateExelStats(id int64, userName string) (string, error) {

	stats := excelize.NewFile()

	defer func() {
		if err := stats.Close(); err != nil {
			slog.Info("Stats file close err:", err)
		}

	}()

	sheetName := fmt.Sprintf("Статистика %s", userName)

	_, err := stats.NewSheet(sheetName)
	if err != nil {
		slog.Error("NewSheet Error:", err)
		return "", err
	}

	trainings, err := u.GetTrainings(id)
	if err != nil {
		slog.Error("GetTrainings Error:", err)

	}

	var date []time.Time
	for i := 0; i < len(trainings); i++ {
		date = append(date, trainings[i].Start)
	}

	start, end := trainingsSort(date)

	earliestTraining := internal.FormatDate(start)
	latestTraining := internal.FormatDate(end)

	averageTrainingLenght, err := u.GetAverageTrainingsLenght(id)
	if err != nil {
		slog.Warn("GetAverageTrainingsLenght Error:", err)
	}

	averageExercisesPerTraining, err := u.GetAverageExercisesPerTraining(id)
	if err != nil {
		slog.Warn("GetAverageExercisesPerTraining Error:", err)
	}

	averageSetsPerTraining, err := u.GetAverageSetsPerTraining(id)
	if err != nil {
		slog.Warn("GetAverageSetsPerTraining Error:", err)
	}

	MostPopularExercise, err := u.GetMostPopularExercise(id)
	if err != nil {
		slog.Warn("GetMostPopularExercise Error:", err)
	}

	LeastPopularExercise, err := u.GetLeastPopularExercise(id)
	if err != nil {
		slog.Warn("GetLeastPopularExercise Error:", err)
	}

	stats.SetColWidth(sheetName, "A", "A", 50)
	stats.SetColWidth(sheetName, "B", "B", 35)
	stats.SetColWidth(sheetName, "C", "C", 42)
	stats.SetColWidth(sheetName, "D", "D", 33)
	stats.SetColWidth(sheetName, "E", "E", 15)

	stats.SetCellValue(sheetName, "A3", fmt.Sprintf("Количество тренрировок за с %s по %s", earliestTraining, latestTraining))
	stats.SetCellValue(sheetName, "B3", len(trainings))
	stats.SetCellValue(sheetName, "A5", "СРЕДНЯЯ ПРОДОЛЖИТЕЛЬНОСТЬ ТРЕНИРОВКИ")
	stats.SetCellValue(sheetName, "B5", averageTrainingLenght)
	stats.SetCellValue(sheetName, "A7", "МАКСИМАЛЬНЫЙ СТРИК")
	stats.SetCellValue(sheetName, "B7", "TODO")
	stats.SetCellValue(sheetName, "A9", "СРЕДНЕЕ КОЛИЧЕСТВО УПРАЖНЕНИЙ ЗА ТРЕНИРОВКУ")
	stats.SetCellValue(sheetName, "B9", averageExercisesPerTraining)
	stats.SetCellValue(sheetName, "A11", "СРЕДНЕЕ КОЛИЧЕСТВО СЭТОВ ЗА ТРЕНИРОВКУ")
	stats.SetCellValue(sheetName, "B11", averageSetsPerTraining)
	stats.SetCellValue(sheetName, "A13", "САМОЕ ПОПУЛЯРНОЕ УПРАЖНЕНИЕ")
	stats.SetCellValue(sheetName, "B13", MostPopularExercise)
	stats.SetCellValue(sheetName, "A15", "САМОЕ НЕПОПУЛЯРНОЕ УПРАЖНЕНИЕ ")
	stats.SetCellValue(sheetName, "B15", LeastPopularExercise)
	stats.SetCellValue(sheetName, "A17", "СТАТИСТИКА ПО КАЖДОМУ УПРАЖНЕНИЮ")
	stats.SetCellValue(sheetName, "A18", "НАЗВАНИЕ УПРАЖНЕНИЯ")
	stats.SetCellValue(sheetName, "B18", "СЭТОВ БЫЛО СДЕЛАНО ЗА ВСЕ ВРЕМЯ")
	stats.SetCellValue(sheetName, "C18", "СРЕДНЕЕ КОЛИЧЕСТВО СЭТОВ ЗА ТРЕНИРОВКУ")
	stats.SetCellValue(sheetName, "D18", "СРЕДНЕЕ КОЛИЧЕСТВО ПОВТОРЕНИЙ")
	stats.SetCellValue(sheetName, "E18", "СРЕДНИЙ ВЕС")

	exercices, err := u.GetExercises(id)
	if err != nil {
		slog.Warn("GetExercises Error:", err)
	}

	for i := 0; i < len(exercices); i++ {

		row := 19 + i

		totalSets, err := u.GetTotalSetsPerExercise(id, exercices[i])
		if err != nil {
			slog.Warn("GetTotalSetsPerExercise error in statsBuilder:", err)
		}

		avgSets, err := u.GetAverageSetsPerExerise(id, exercices[i])
		if err != nil {
			slog.Warn("getAvgSetsPerExercise error in statsBuilder", err)
		}

		avgReps, err := u.GetAverageReps(id, exercices[i])
		if err != nil {
			slog.Warn("GetAverageReps Error:", err)
		}

		avgWeight, err := u.GetAverageWeight(id, exercices[i])
		if err != nil {
			slog.Warn("GetAverageWeight Error:", err)
		}

		stats.SetCellValue(sheetName, fmt.Sprintf("A%d", row), exercices[i])
		stats.SetCellValue(sheetName, fmt.Sprintf("B%d", row), totalSets)
		stats.SetCellValue(sheetName, fmt.Sprintf("C%d", row), avgSets)
		stats.SetCellValue(sheetName, fmt.Sprintf("D%d", row), avgReps)
		stats.SetCellValue(sheetName, fmt.Sprintf("E%d", row), fmt.Sprintf("%.2f", avgWeight))
	}

	filePath := fmt.Sprintf("%s_stats.xlsx", userName)

	if err := stats.SaveAs(filePath); err != nil {
		slog.Error("SaveAs Error:", err)
		return "", err
	}

	return filePath, err

}

func trainingsSort(trainingsStart []time.Time) (first time.Time, last time.Time) {

	sort.Slice(trainingsStart, func(i, j int) bool {
		return trainingsStart[i].Before(trainingsStart[j])

	})

	return trainingsStart[0], trainingsStart[len(trainingsStart)-1]

}
