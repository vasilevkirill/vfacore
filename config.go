package vfacore

import (
	"github.com/spf13/viper"
	"log"
)

var (
	configGlobalS configS
	configViper   = viper.New()
)

func loadConfig() error {
	Config := viper.New()
	Config.AddConfigPath("./config")
	Config.SetConfigName("config")
	Config.SetConfigType("yaml")

	err := Config.ReadInConfig()
	if err != nil {
		return err
		//panic(fmt.Errorf("fatal error config file: %s", err.Error()))
	}
	var configSv configS
	err = Config.Unmarshal(&configSv)
	if err != nil {
		return err
	}
	configGlobalS = configSv
	configViper = Config
	if err = checkConfigLdap(); err != nil {
		return err
	}
	if err = checkConfigTelegram(); err != nil {
		return err
	}
	return nil
}

func checkConfigLdap() error {

	if configGlobalS.Ldap.Dn == "" {
		return errorGetFromId(400)
	}
	if configGlobalS.Ldap.User == "" {
		return errorGetFromId(401)
	}
	if configGlobalS.Ldap.Password == "" {
		return errorGetFromId(402)
	}

	return nil
}

func checkConfigTelegram() error {
	if configGlobalS.Telegram.Token == "" {
		return errorGetFromId(303)
	}
	if configGlobalS.Telegram.HookPort == 0 {
		return errorGetFromId(304)
	}
	if configGlobalS.Telegram.HookPort > 65535 && configGlobalS.Telegram.HookPort < 1 {
		return errorGetFromId(305)
	}
	if configGlobalS.Telegram.HookDomain == "" {
		return errorGetFromId(306)
	}
	if configGlobalS.Telegram.PoolAddress == "" {
		log.Println("telegram.PoolAddress config is empty, set default 0.0.0.0")
		configGlobalS.Telegram.PoolAddress = "0.0.0.0"
	}
	if configGlobalS.Telegram.HookCertPub == "" || configGlobalS.Telegram.HookCertKey == "" {
		log.Println("telegram.HookCertPub or telegram.HookCertKey config is empty, used self-sign")
		err := generateCertificate()
		if err != nil {
			return errorGetFromIdAddSuffix(500, err.Error())
		}
	}
	return nil
}
