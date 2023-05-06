package vfacore

import (
	"crypto/tls"
	"errors"
	"fmt"
	ldapv3 "github.com/go-ldap/ldap/v3"
	"github.com/robfig/cron/v3"
	"log"
	"strconv"
)

var (
	conn *ldapv3.Conn
)

var ()

func ldapRun() error {
	err := createConn()
	if err != nil {
		return err
	}
	return nil
}
func ping() {
	filer := "(&(objectClass=organizationalPerson)(objectClass=user)(sAMAccountName=test))"
	attrs := []string{"displayName"}
	_, err := searchFilterAttrs(filer, attrs)
	if err != nil {
		log.Println("Ошибка ping ldap пытаемся установить новое соединение с ldap")
		err := createConn()
		log.Println(err)
		return
	}
	return
}

func createConn() error {
	for _, ldapServer := range configGlobalS.Ldap.Servers {
		log.Printf("Ldap - Попытка установить соединение с сервером ldap: %s", ldapServer)
		ldapL, err := ldapv3.Dial("tcp", ldapServer)
		if err != nil {
			log.Printf("Ldap - Ошибка при подключении к серверу ldap: %s", ldapServer)
			log.Println(err)
			continue
		}
		err = ldapL.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			log.Printf("Ldap - Ошибка при установки TLS соединения с сервером ldap: %s", ldapServer)
			log.Println(err)
			continue
		}
		err = ldapL.Bind(configGlobalS.Ldap.User, configGlobalS.Ldap.Password)
		if err != nil {
			log.Printf("Ldap - Ошибка при авторизиции (bind) на сервере ldap: %s", ldapServer)
			log.Println(err)
			continue
		}
		log.Printf("Ldap - Установленно соединение с сервером ldap: %s", ldapServer)
		conn = ldapL
		schedule()
		return nil
	}

	return errorGetFromId(801)
}
func schedule() {
	c := cron.New()
	_, _ = c.AddFunc("@every 1m", ping)
	c.Start()
}

func searchFilterAttrs(filter string, attr []string) ([]*ldapv3.Entry, error) {
	searchRequest := ldapv3.NewSearchRequest(
		configGlobalS.Ldap.Dn,
		ldapv3.ScopeWholeSubtree,
		ldapv3.NeverDerefAliases,
		0,
		0,
		false,
		filter,
		attr,
		nil,
	)
	req, err := conn.Search(searchRequest)
	if err != nil {
		errN := errorGetFromIdAddSuffix(802, err.Error())
		return []*ldapv3.Entry{}, errN
	}
	return req.Entries, nil
}

func (u *ldapUser) PullViaTelegramId() error {
	if u.TelegramId == 0 {
		return errors.New("telegramid не указан")
	}
	filer := fmt.Sprintf("(&(objectClass=organizationalPerson)(objectClass=user)(pager=%d))", u.TelegramId)
	attrs := []string{"pager", "displayName", "sAMAccountName"}
	req, err := searchFilterAttrs(filer, attrs)
	if err != nil {
		return err
	}
	if len(req) == 0 {
		return ldapErrUserNotFound
	}
	if len(req) > 1 {
		return ldapErrUserFoundMoreThanOne
	}
	u.DispalyName = req[0].GetAttributeValue("displayName")
	u.SAMAccountName = req[0].GetAttributeValue("sAMAccountName")
	return nil
}

func (u *ldapUser) PullViaSAMAccountName() error {
	if u.SAMAccountName == "" {
		return ldapErrUserSAMAccountRequired
	}
	filer := fmt.Sprintf("(&(objectClass=organizationalPerson)(objectClass=user)(sAMAccountName=%s))", u.SAMAccountName)
	attrs := []string{"pager", "displayName", "sAMAccountName"}
	req, err := searchFilterAttrs(filer, attrs)
	if err != nil {
		return err
	}
	if len(req) == 0 {
		return ldapErrUserNotFound
	}
	if len(req) > 1 {
		return ldapErrUserFoundMoreThanOne
	}
	u.DispalyName = req[0].GetAttributeValue("displayName")
	pager := req[0].GetAttributeValue("pager")
	TelegramId, err := strconv.ParseInt(pager, 10, 64)
	if err != nil {
		errN := errorGetFromIdAddSuffix(803, err.Error())
		return errN
	}
	u.TelegramId = TelegramId
	return nil
}
