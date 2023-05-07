package vfacore

func Run() error {
	err := loadConfig() // поулчаем конфиг
	if err != nil {
		return err
	}

	initQ()

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
