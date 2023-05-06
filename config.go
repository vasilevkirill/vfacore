package _fa

import (
	"core/certificate"
	"core/errc"
	"core/ldap"
	"core/tg"
	"fmt"
	"github.com/spf13/viper"
	"log"
)

func loadConfig() *viper.Viper {
	Config := viper.New()
	Config.AddConfigPath("./config")
	Config.SetConfigName("config")
	Config.SetConfigType("yaml")

	err := Config.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error config file: %s", err.Error()))
	}
	return Config
}

func checkConfigLdap(config *ldap.Config) error {

	if config.Dn == "" {
		return errc.GetError(400)
	}
	if config.User == "" {
		return errc.GetError(401)
	}
	if config.Password == "" {
		return errc.GetError(402)
	}

	return nil
}

func checkConfigTelegram(config *tg.Config) error {
	if config.Token == "" {
		return errc.GetError(303)
	}
	if config.HookPort == 0 {
		return errc.GetError(304)
	}
	if config.HookPort > 65535 && config.HookPort < 1 {
		return errc.GetError(305)
	}
	if config.HookDomain == "" {
		return errc.GetError(306)
	}
	if config.PoolAddress == "" {
		log.Println("telegram.PoolAddress config is empty, set default 0.0.0.0")
		config.PoolAddress = "0.0.0.0"
	}
	if config.HookCertPub == "" || config.HookCertKey == "" {
		log.Println("telegram.HookCertPub or telegram.HookCertKey config is empty, user self-sign")
		err := certificate.GenerateCertificate(config.HookDomain)
		if err != nil {
			return errc.ErrorAddSuffix(500, err.Error())
		}
	}
	return nil
}
