package vfacore

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"log"
	"net/http"
)

var bot *tgbotapi.BotAPI

func telegramRun() error {

	bt, err := tgbotapi.NewBotAPI(configGlobalS.Telegram.Token)
	if err != nil {
		return errorGetFromIdAddSuffix(600, err.Error())

	}
	bt.Debug = configGlobalS.Telegram.Debug
	webHookAddress := fmt.Sprintf("https://%s:%d", configGlobalS.Telegram.HookDomain, configGlobalS.Telegram.HookPort)
	configGlobalS.Telegram.WebHookAddress = webHookAddress
	wh, err := tgbotapi.NewWebhookWithCert(webHookAddress, tgbotapi.FilePath(configGlobalS.Telegram.HookCertPub))
	if err != nil {
		return errorGetFromIdAddSuffix(601, err.Error())
	}
	_, err = bt.Request(wh)
	if err != nil {
		return errorGetFromIdAddSuffix(602, err.Error())

	}
	info, err := bt.GetWebhookInfo()
	if err != nil {
		return errorGetFromIdAddSuffix(603, err.Error())
	}

	if info.LastErrorDate != 0 {
		log.Printf("Telegram callback failed: %s", info.LastErrorMessage)
	}
	bot = bt
	go updatesWord()

	return nil
}

func updatesWord() {
	updates := bot.ListenForWebhook("/")

	go runHttpServer()
	for update := range updates {
		if checkCallbackQuery(update) {
			continue
		}
		if update.Message == nil { // ignore any non-Message updates
			continue
		}
		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}
		if update.Message.From.IsBot { //ignore bot
			continue
		}
		switch update.Message.Command() {
		case "start":
			cmdStart(update)
		default:
			debug(fmt.Sprintf("Получена комманда %s", update.Message.Command()))
		}
	}
}
func checkCallbackQuery(update tgbotapi.Update) bool {
	CallbackQuery := update.CallbackQuery
	data := ""
	if CallbackQuery != nil {
		data = CallbackQuery.Data
	}

	if data == "" {
		return false
	}
	msg := CallbackQuery.Message
	debug(fmt.Sprintf("Пользователь %d нажал %s", msg.Chat.ID, data))
	err := removeMsg(msg)
	if err != nil {
		log.Println(err)
		return true
	}
	m := qu.GetMsg(msg.Chat.ID)
	if data == "yes" {
		m.Chan <- 1
		return true
	}

	if data == "no" {
		m.Chan <- 0
		return true
	}
	return false
}

func sendQuery(user ldapUser, timeout int) error {

	inlineKeyboard := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Да", "yes"),
			tgbotapi.NewInlineKeyboardButtonData("Нет", "no"),
		),
	)
	str := fmt.Sprintf("Кто-то пытается авторизоваться под вашей учетной записью\nЭто вы?\n Необходимо ответить в течении %d секунд", timeout)
	msg := tgbotapi.NewMessage(user.TelegramId, str)
	msg.ReplyMarkup = inlineKeyboard
	msgSend, err := bot.Send(msg)
	if err != nil {
		return err
	}
	qu.SetMsgId(msgSend.Chat.ID, int64(msgSend.MessageID))
	return nil
}

func removeMsg(msg *tgbotapi.Message) error {
	deleteMsgConfig := tgbotapi.NewDeleteMessage(msg.Chat.ID, msg.MessageID)
	_, err := bot.Request(deleteMsgConfig)
	if err != nil {

		return errorGetFromIdAddSuffix(604, err.Error())
	}
	return nil
}
func removeMsgByChaiIDMsgIDForce(chatid, msgid int64) {
	deleteMsgConfig := tgbotapi.NewDeleteMessage(chatid, int(msgid))
	_, _ = bot.Request(deleteMsgConfig)
	return
}

func runHttpServer() {
	strConnect := fmt.Sprintf("%s:%d", configGlobalS.Telegram.PoolAddress, configGlobalS.Telegram.PoolPort)
	err := http.ListenAndServeTLS(strConnect, configGlobalS.Telegram.HookCertPub, configGlobalS.Telegram.HookCertKey, nil)
	if err != nil {
		errN := errorGetFromIdAddSuffix(605, err.Error(), strConnect)
		log.Panic(errN)
	}
}

func debug(str string) {
	if configGlobalS.Telegram.Debug {
		log.Println(str)
	}
}

func cmdStart(update tgbotapi.Update) {

	debug("Система получила команду /start")

	msgWait := tgbotapi.NewMessage(update.Message.Chat.ID, "Ждите...")
	msgW, err := bot.Send(msgWait)

	user := ldapUser{}
	user.TelegramId = update.Message.From.ID
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
	err = user.PullViaTelegramId()
	if err != nil {
		switch err {
		case ldapErrUserNotFound:
			msg.Text = fmt.Sprintf("Привет, мы не знакомы\n отправь с службу поддержки твой ID %d ", user.TelegramId)
		case ldapErrUserFoundMoreThanOne:
			msg.Text = fmt.Sprintf("Сообщение об ошибке с кодом 100, пожалуйста, обратитесь в службу технической поддержки.\nid:%d", user.TelegramId)
		default:
			msg.Text = fmt.Sprintf("Сообщение об ошибке с кодом 999, пожалуйста, обратитесь в службу технической поддержки.\nid:%d", user.TelegramId)
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
	msg.Text = "Здравствуйте, всё подготовлено, мы уже знакомы."
	_, err = bot.Send(msg)
	if err != nil {
		log.Println(err)
		return
	}

}
