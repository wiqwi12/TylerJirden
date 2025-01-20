package telegram

import (
	"GymBot/internal/application"
	"fmt"
	"gopkg.in/telebot.v3"
	"log/slog"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	weightRegex = `^([1-9]\d*|0)(\.\d+)?([eE][+-]?\d+)?$`
	repsRegex   = `^(0|[1-9]\d*)$`
)

var (
	repsRegexp   = regexp.MustCompile(repsRegex)
	weightRegexp = regexp.MustCompile(weightRegex)
)

type BotHandler struct {
	Service *application.Service
}

func NewBotHandler(service *application.Service) *BotHandler {
	return &BotHandler{
		Service: service,
	}
}

func (b *BotHandler) MsgMainHandler(c telebot.Context) error {
	msg := c.Message().Text

	switch msg {
	case "/start":

		b.StartHandler(c)
	default:

		c.Send("Неизвестная команда", StartKeyboard())
	}

	return nil
}

func (b *BotHandler) DataHandler(c telebot.Context) error {

	var err error
	data := c.Data()

	switch {
	case strings.HasPrefix(data, "next_"):

		page, err := strconv.Atoi(strings.TrimPrefix(data, "next_"))
		if err != nil {
			slog.Error("strconv err:", err)
		}

		c.Edit("Выберите упражнение", b.PagKeyboard(c.Sender().ID, int64(page)))

	case strings.HasPrefix(data, "prev_"):

		page, err := strconv.Atoi(strings.TrimPrefix(data, "prev_"))
		if err != nil {
			slog.Error("strconv err:", err)
		}

		c.Edit("Выберите упражнение", b.PagKeyboard(c.Sender().ID, int64(page)))

	case strings.HasPrefix(data, "exercise_"):

		exercise := strings.TrimPrefix(data, "exercise_")

		err := b.Service.Repo.SetExercise(c.Sender().ID, exercise)
		if err != nil {
			slog.Error("Set exercise:", err)
		}

		c.Edit("Упражнение выбрано! Можете начинать!", TrainingKeyboardWithExerciseChosen())

	default:
		switch data {
		case "start_training":
			err = b.StartTrainingHandler(c)
		case "add_exercise":
			c.Send("Введите упражнение")
			c.Bot().Handle(telebot.OnText, b.AddExerciseHandler)
		case "end_training":
			err = b.EndTrainingHandler(c)
		case "start_set":
			err = b.StartSetHandler(c)
		case "end_set":
			err = b.EndSetHandler(c)
		case "choose_exercise":
			err = c.Edit("Выберите упражнение", b.PagKeyboard(c.Sender().ID, 1))
		case "show_stats":
			err = b.StatsHandler(c)

		}
	}

	if err != nil {
		slog.Error("Error in Data Handler:", err)
		return err
	}

	return nil
}

func (b *BotHandler) StartHandler(c telebot.Context) error {

	Exsist, err := b.Service.Repo.UserCheck(c.Sender().ID)
	if err != nil {
		slog.Error("User check error:", err)
	}

	if !Exsist {
		err = b.Service.Repo.RegisterUser(c.Sender().ID)

		if err != nil {
			slog.Error("User registration err:", err)
			return err
		}

		c.Send("Привет! Бот позволяет тебе добавлять упраженияи, трекать тренировку и получать статистику", StartKeyboard())

	} else if Exsist {
		c.Send("Давно не виделись!", StartKeyboard())
	} else {
		c.Send("Неизвестная команда", StartKeyboard())
	}
	return nil
}

func (b *BotHandler) StartTrainingHandler(c telebot.Context) error {

	err := b.Service.Repo.StartTrainig(c.Sender().ID, c.Message().Time())
	if err != nil {
		slog.Error("Start training error:", err)
		return err
	}

	c.Edit("Тренировка началась!", TrainingKeyboard())

	return nil

}

func (b *BotHandler) EndTrainingHandler(c telebot.Context) error {

	err := b.Service.Repo.EndTraining(c.Sender().ID, c.Message().Time())
	if err != nil {
		slog.Error("End training error:", err)
		return err
	}

	c.Edit("Тренировка завершена!", StartKeyboard())

	return nil
}

func (b *BotHandler) StartSetHandler(c telebot.Context) error {

	isChosen, err := b.Service.Repo.IsExerciseChoosen(c.Sender().ID)
	if err != nil {
		slog.Error("Is exercise choosen error:", err)
		return err
	}

	if !isChosen {
		c.Edit("Для начала сэта выберите упражнение.", ChooseKeyboard())
	} else if isChosen {
		err = b.Service.Repo.StartSet(c.Sender().ID, c.Message().Time())
		if err != nil {
			slog.Error("start set error", err)
		}
		c.Edit("Сэт идет!", SetKeyboard())

	}

	return nil
}

func (b *BotHandler) EndSetHandler(c telebot.Context) error {

	b.Service.Repo.EndSet(c.Sender().ID, c.Message().Time())

	c.Send("Сэт завершен! Введите вес, который вы использовали.")

	c.Bot().Handle(telebot.OnText, b.WeightHandler)

	return nil
}

func (b *BotHandler) WeightHandler(c telebot.Context) error {

	msg := strings.ReplaceAll(c.Message().Text, ",", ".")

	if weightRegexp.MatchString(msg) {

		weight, err := strconv.ParseFloat(msg, 64)
		if err != nil {
			slog.Error("parse float error:", err)
			c.Send("Ошибка ввода веса. Пожалуйста, введите число c одной цифрой после запятой(точка тож сойдет).")
			c.Bot().Handle(telebot.OnText, b.WeightHandler)
		}

		err = b.Service.Repo.SetWeight(c.Sender().ID, weight)
		if err != nil {
			slog.Error("set weight err:", err)
			return err
		}

		c.Send("Теперь введите количество пfовторений.")
		c.Bot().Handle(telebot.OnText, b.RepsHandler)

	} else if !weightRegexp.MatchString(msg) {

		c.Send("Ошибка ввода веса. Пожалуйста, введите число c одной цифрой после запятой(точка тож сойдет).")

		c.Bot().Handle(telebot.OnText, b.WeightHandler)
	}

	return nil
}

func (b *BotHandler) RepsHandler(c telebot.Context) error {

	if repsRegexp.MatchString(c.Message().Text) {

		reps, err := strconv.Atoi(c.Message().Text)
		if err != nil {
			c.Send("Ошибка ввода повторений. Пожалуйста, введите целое число.")
			c.Bot().Handle(telebot.OnText, b.RepsHandler)
		}

		b.Service.Repo.SetReps(c.Sender().ID, reps)
		c.Send("Сэт успешно завершен! Все данные затреканы!", TrainingKeyboard())

	} else if !repsRegexp.MatchString(c.Message().Text) {

		c.Send("Ошибка ввода повторений. Пожалуйста, введите целое число.")
		c.Bot().Handle(telebot.OnText, b.RepsHandler)

	}

	return nil

}

func (b *BotHandler) AddExerciseHandler(c telebot.Context) error {

	err := b.Service.Repo.AddExercise(c.Sender().ID, c.Message().Text)
	if err != nil {
		slog.Error("add exercise error:", err)
		return err
	}

	isActive, err := b.Service.Repo.IsTrainingActive(c.Sender().ID)
	if err != nil {
		slog.Error("add exercise error:", err)
		return err
	}

	if isActive {
		c.Send(fmt.Sprintf("Упражнение '%s' добавлено.", c.Message().Text), TrainingKeyboard())
	} else {
		c.Send(fmt.Sprintf("Упражнение '%s' добавлено.", c.Message().Text), StartKeyboard())
	}

	return nil

}

func (b *BotHandler) StatsHandler(c telebot.Context) error {

	filePath, err := b.Service.Repo.GenerateExelStats(c.Sender().ID, c.Sender().Username)
	if err != nil {
		slog.Error("generate exel stats error:", err)
		return err
	}

	file, err := os.Open(filePath)
	if err != nil {
		slog.Error("open exel stats error:", err)
		return err
	}

	defer file.Close()

	doc := &telebot.Document{
		File:     telebot.FromReader(file),
		FileName: filePath,
	}

	c.Send(doc)

	return nil

}
