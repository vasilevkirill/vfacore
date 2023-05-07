package vfacore

func Run(shb bool) error {
	err := loadConfig() // поулчаем конфиг
	if err != nil {
		return err
	}

	initQ()
	if shb {
		go showBanner()
	}
	err = ldapRun() // инициалзируем ldap
	if err != nil {
		return err
	}

	err = telegramRun() // инициалзируем бета телеграм
	if err != nil {
		return err
	}

	err = radiusRun() // инициалзируем радиус сервер
	if err != nil {
		return err
	}

	return nil
}
