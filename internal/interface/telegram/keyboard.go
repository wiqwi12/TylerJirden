package telegram

import (
	"fmt"
	"gopkg.in/telebot.v3"
	"log/slog"
)

var (
	btnStartTraining = telebot.InlineButton{
		Text: "Начать тренировку",
		Data: "start_training",
	}

	btnEndTraining = telebot.InlineButton{
		Text: "Закончить тренировку",
		Data: "end_training",
	}

	btnStartSet = telebot.InlineButton{
		Text: "Начать сэт",
		Data: "start_set",
	}

	btnEndSet = telebot.InlineButton{
		Text: "Закончить сэт",
		Data: "end_set",
	}

	btnAdd = telebot.InlineButton{
		Text: "Добавить упражнение",
		Data: "add_exercise",
	}

	btnChooseExercise = telebot.InlineButton{
		Text: "Выбрать упражнение",
		Data: "choose_exercise",
	}
)

func StartKeyboard() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{btnStartTraining}, {btnAdd},
		}}
}

func TrainingKeyboard() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{btnStartSet, btnEndTraining},
			{btnAdd, btnChooseExercise},
		}}
}

func TrainingKeyboardWithExerciseChosen() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{btnStartSet, btnEndTraining},
			{btnAdd},
		}}
}

func SetKeyboard() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{btnEndSet},
		}}
}

func ChooseKeyboard() *telebot.ReplyMarkup {
	return &telebot.ReplyMarkup{
		InlineKeyboard: [][]telebot.InlineButton{
			{btnChooseExercise},
		},
	}
}

func (b *BotHandler) PagKeyboard(id, current_page int64) *telebot.ReplyMarkup {

	var (
		nexBtn = telebot.InlineButton{
			Text: "Next",
			Data: fmt.Sprintf("next_%d", current_page+1),
		}

		prevBtn = telebot.InlineButton{
			Text: "Previous",
			Data: fmt.Sprintf("prev_%d", current_page-1),
		}
	)

	page, err := b.Service.Repo.GetPage(id, current_page)
	if err != nil {
		slog.Error("GetPage err:", err)
	}

	var exerciseBtns []telebot.InlineButton
	for _, exercise := range page {
		exerciseBtns = append(exerciseBtns, telebot.InlineButton{
			Text: exercise,
			Data: fmt.Sprintf("exercise_%s", exercise),
		})
	}

	rows := [][]telebot.InlineButton{}
	for _, btn := range exerciseBtns {
		rows = append(rows, []telebot.InlineButton{btn})
	}

	maxPage, err := b.Service.Repo.MaxPages(id)
	if err != nil {
		slog.Error("GetMaxPages err:", err)
	}

	if current_page == 1 && current_page != maxPage {
		rows = append(rows, []telebot.InlineButton{nexBtn})
	} else if current_page == maxPage && current_page != 1 {
		rows = append(rows, []telebot.InlineButton{prevBtn})
	} else if current_page == maxPage && current_page == 1 {
		slog.Info("Keyboard rows generated", slog.Any("Rows", rows))
		return &telebot.ReplyMarkup{
			InlineKeyboard: rows,
		}
	} else {
		rows = append(rows, []telebot.InlineButton{prevBtn, nexBtn})
	}

	return &telebot.ReplyMarkup{
		InlineKeyboard: rows,
	}
}
