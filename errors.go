package vfacore

import (
	"errors"
	"fmt"
)

var (
	errorMap = map[int]string{
		300: "private key Null string. show in config telegram.HookCertKey",
		301: "public key Null string. show in config telegram.HookCertPub",
		302: "domain Null string. show in config telegram.HookDomain",
		303: "parameter in telegram.token config is empty",
		304: "parameter in telegram.HookPort config is empty or NUll",
		305: "parameter in telegram.HookPort config min 1 max 65535",
		306: "parameter in telegram.HookDomain config is empty",
		400: "parameter in ldap.dn config is empty",
		401: "parameter in ldap.user config is empty",
		402: "parameter in ldap.password config is empty",
		500: "Generate Certificate:",
	}
	ldapErrUserNotFound           = errors.New("Пользователь не найден")
	ldapErrUserFoundMoreThanOne   = errors.New("Найденно пользователей больше чем один")
	ldapErrUserSAMAccountRequired = errors.New("SAMAccount обезателен")
)

func errorGetFromId(i int) error {
	e, ok := errorMap[i]
	if !ok {
		return errors.New("unknown error")
	}
	er := fmt.Sprintf("error:#%d > %s", i, e)
	return errors.New(er)
}

func errorGetFromIdAddSuffix(i int, suffix string) error {

	err := errorGetFromId(i)
	errNew := fmt.Sprintf("%s %s", err.Error(), suffix)
	return errors.New(errNew)
}
