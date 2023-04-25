package tg

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"main/ldap"
)

func cmdStart(update tgbotapi.Update) {

	debug("Получена комманда /start")

	msgWait := tgbotapi.NewMessage(update.Message.Chat.ID, "Ждите...")
	msgW, err := bot.Send(msgWait)

	user := ldap.User{}
	user.TelegramId = update.Message.From.ID
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	err = user.PullViaTelegramId()
	if err != nil {
		switch err {
		case ldap.ErrUserNotFound:
			msg.Text = fmt.Sprintf("Привет, мы не знакомы\n отправь с службу поддержки #100 твой ID %d ", user.TelegramId)
		case ldap.ErrUserFoundMoreThanOne:
			msg.Text = fmt.Sprintf("Произошла ошибка 100, свяжись с технической поддержкой\nid:%d", user.TelegramId)
		default:
			msg.Text = fmt.Sprintf("Произошла ошибка 999, свяжись с технической поддержкой\nid:%d", user.TelegramId)
		}
		err = removeMsg(&msgW)
		if err != nil {
			log.Println(err)
			return
		}
		_, err = bot.Send(msg)
		if err != nil {
			log.Println(err)
			return
		}
		return
	}
	err = removeMsg(&msgW)
	if err != nil {
		log.Println(err)
	}
	msg.Text = "Привет всё готовы, мы уже знакомы"
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}

}
