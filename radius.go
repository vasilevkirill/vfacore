package vfacore

import (
	"context"
	"errors"
	"fmt"
	"layeh.com/radius"
	"layeh.com/radius/rfc2865"
	"log"
	"strings"
	"time"
)

func radiusRun() error {

	server := radius.PacketServer{
		Handler:      radius.HandlerFunc(handler),
		SecretSource: radius.StaticSecretSource([]byte(configGlobalS.Radius.Secret)),
		Addr:         fmt.Sprintf("%s:%d", configGlobalS.Radius.Address, configGlobalS.Radius.Port),
	}

	log.Printf(fmt.Sprintf("Запуск сервера на %s", server.Addr))
	if err := server.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

func handler(w radius.ResponseWriter, r *radius.Request) {

	username := rfc2865.UserName_GetString(r.Packet)
	username = getUserName(username)
	user, err := getUser(username)
	if err != nil {
		log.Println(err)
		sendAccessReject(w, r)
		return
	}
	if user.TelegramId == 0 {
		log.Printf("Пользователь %s не имеет значения TelegramId", user.SAMAccountName)
		sendAccessReject(w, r)
		return
	}
	if qu.IssetKey(user.TelegramId) {
		log.Printf("Запрос пользователю %s уже отправлен", user.SAMAccountName)
		return
	}
	log.Printf("Запрос на подключение от пользователя %s", user.SAMAccountName)
	qu.AddKey(user.TelegramId)

	err = sendQuery(user, configGlobalS.Radius.Answertimeout)
	if err != nil {
		log.Println(err)
		sendAccessReject(w, r)
		return
	}

	ctx := context.Background()
	ctx, cancelFunctionContext := context.WithTimeout(ctx, time.Duration(configGlobalS.Radius.Answertimeout)*time.Second)
	defer func() {
		qu.RemoveKey(user.TelegramId)
		cancelFunctionContext()
		sendAccessReject(w, r)
		return
	}()
	msg := qu.GetMsg(user.TelegramId)
	err = waitAnswer(ctx, msg, user)
	if err != nil {
		qu.RemoveKey(user.TelegramId)
		log.Println(err)
		sendAccessReject(w, r)
		removeMsgByChaiIDMsgIDForce(user.TelegramId, msg.MsgId)
		return
	}
	qu.RemoveKey(user.TelegramId)
	log.Printf("Пользователь %s aвторизирован", user.SAMAccountName)
	sendAccessAccept(w, r)
}
func waitAnswer(ctx context.Context, msg queueMsg, user ldapUser) error {
	for {
		select {
		case <-ctx.Done():
			//timeout
			return errors.New(fmt.Sprintf("Пользователю %s отказанно: %s", user.SAMAccountName, ctx.Err()))

		case num := <-msg.Chan:
			if num == 0 {
				return errors.New(fmt.Sprintf("Пользователь %s выбрал No", user.SAMAccountName))
			}
			return nil
		}
	}
}
func sendAccessAccept(w radius.ResponseWriter, r *radius.Request) {
	send(w, r, radius.CodeAccessAccept)
	return
}

func sendAccessReject(w radius.ResponseWriter, r *radius.Request) {
	send(w, r, radius.CodeAccessReject)
	return
}

func send(w radius.ResponseWriter, r *radius.Request, code radius.Code) {
	p := r.Response(code)
	prx := rfc2865.ProxyState_Get(r.Packet)
	p.Add(rfc2865.ProxyState_Type, prx)
	err := w.Write(p)
	if err != nil {
		log.Printf("Radius send error: %s", err.Error())
	}
}

func getUserName(user string) string {
	userSplit := strings.Split(user, `\`)
	if len(userSplit) == 2 {
		return userSplit[1]
	}
	return user
}

func getUser(sAMAccountName string) (ldapUser, error) {
	u := ldapUser{}
	u.SAMAccountName = sAMAccountName
	err := u.PullViaSAMAccountName()
	if err != nil {
		return ldapUser{}, err
	}
	return u, nil
}
