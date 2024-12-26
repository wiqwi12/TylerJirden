package postgres

import (
	"GymBot/internal/domain/repository"
	"database/sql"
	"log"
	"log/slog"
	"time"

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
	}

	_, err = u.Db.Query(query, args...)
	if err != nil {
		slog.Error("Start Trainig Query Error:", err)
	}

	return err
}

func (u *UserRepositoryDB) EndTraining(id int64, endTime time.Time) error {
	slog.Info("End training started")

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
	slog.Info("Start set started")

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

	slog.Info("QUERY:", query, "args:", args)

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

	increment := squirrel.Select("MAX(exercise_id)").From("exercises").Where(squirrel.Eq{"user_id": id}).PlaceholderFormat(squirrel.Dollar)

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
