package _fa

import (
	ldap2 "core/ldap"
	"core/queue"
	radius2 "core/radius"
	tg2 "core/tg"
)

func Run() error {
	config := loadConfig() // поулчаем конфиг
	var LdapConfig ldap2.Config
	err := config.UnmarshalKey("ldap", &LdapConfig)
	if err != nil {
		return err
	}
	err = checkConfigLdap(&LdapConfig)
	if err != nil {
		return err
	}

	var tgConfig tg2.Config
	err = config.UnmarshalKey("telegram", &tgConfig)
	if err != nil {
		return err
	}
	err = checkConfigTelegram(&tgConfig)
	if err != nil {
		return err
	}

	var radiusConfig radius2.Config
	err = config.UnmarshalKey("radius", &radiusConfig)
	if err != nil {
		return err
	}

	queue.InitQ()

	err = ldap2.Run(LdapConfig) // инициалзируем ldap
	if err != nil {
		return err
	}

	err = tg2.Run(tgConfig) // инициалзируем бета телеграм
	if err != nil {
		return err
	}

	err = radius2.Run(radiusConfig) // инициалзируем радиус сервер
	if err != nil {
		return err
	}
	return nil
}
