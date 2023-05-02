package ldap

import (
	"crypto/tls"
	"errors"
	ldapv3 "github.com/go-ldap/ldap/v3"
	"github.com/robfig/cron/v3"
	"log"
)

var (
	conn   *ldapv3.Conn
	config Config
)

func Run(cfg Config) error {
	config = cfg
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
	for _, ldapServer := range config.Servers {
		log.Printf("Попытка установить соединение с сервером ldap: %s", ldapServer)
		ldapL, err := ldapv3.Dial("tcp", ldapServer)
		if err != nil {
			log.Printf("Ошибка при подключении к серверу ldap: %s", ldapServer)
			log.Println(err)
			continue
		}
		err = ldapL.StartTLS(&tls.Config{InsecureSkipVerify: true})
		if err != nil {
			log.Printf("Ошибка при установки TLS соединения с сервером ldap: %s", ldapServer)
			log.Println(err)
			continue
		}
		err = ldapL.Bind(config.User, config.Password)
		if err != nil {
			log.Printf("Ошибка при авторизиции (bind) на сервере ldap: %s", ldapServer)
			log.Println(err)
			continue
		}
		log.Printf("Установленно соединение с сервером ldap: %s", ldapServer)
		conn = ldapL
		schedule()
		return nil
	}

	return errors.New("не смогли установить соединение с ldap серверами")
}
func schedule() {
	c := cron.New()
	_, _ = c.AddFunc("@every 1m", ping)
	c.Start()
}

func searchFilterAttrs(filter string, attr []string) ([]*ldapv3.Entry, error) {
	searchRequest := ldapv3.NewSearchRequest(
		config.Dn,
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
		return []*ldapv3.Entry{}, err
	}
	return req.Entries, nil
}
