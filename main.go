package main

import (
	"log"
	"main/ldap"
	"main/queue"
	radius "main/radius"
	tg "main/tg"
)

func main() {
	config := LoadConfig() // поулчаем конфиг
	queue.InitQ()
	var LdapConfig ldap.Config
	err := config.UnmarshalKey("ldap", &LdapConfig)
	err = ldap.Run(LdapConfig) // инициалзируем ldap
	checkFatal(err)
	var tgConfig tg.Config
	err = config.UnmarshalKey("telegram", &tgConfig)
	checkFatal(err)
	err = tg.Run(tgConfig) // инициалзируем бета телеграм
	checkFatal(err)

	var radiusConfig radius.Config
	err = config.UnmarshalKey("radius", &radiusConfig)
	checkFatal(err)
	err = radius.Run(radiusConfig) // инициалзируем радиус сервер
	checkFatal(err)

}

func checkFatal(err error) {
	if err == nil {
		return
	}
	log.Panic(err)
}
