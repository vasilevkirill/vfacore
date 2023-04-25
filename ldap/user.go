package ldap

import (
	"errors"
	"fmt"
	"strconv"
)

func (u *User) PullViaTelegramId() error {
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
		return ErrUserNotFound
	}
	if len(req) > 1 {
		return ErrUserFoundMoreThanOne
	}
	u.DispalyName = req[0].GetAttributeValue("displayName")
	u.SAMAccountName = req[0].GetAttributeValue("sAMAccountName")
	return nil
}

func (u *User) PullViaSAMAccountName() error {
	if u.SAMAccountName == "" {
		return ErrUserSAMAccountRequired
	}
	filer := fmt.Sprintf("(&(objectClass=organizationalPerson)(objectClass=user)(sAMAccountName=%s))", u.SAMAccountName)
	attrs := []string{"pager", "displayName", "sAMAccountName"}
	req, err := searchFilterAttrs(filer, attrs)
	if err != nil {
		return err
	}
	if len(req) == 0 {
		return ErrUserNotFound
	}
	if len(req) > 1 {
		return ErrUserFoundMoreThanOne
	}
	u.DispalyName = req[0].GetAttributeValue("displayName")
	pager := req[0].GetAttributeValue("pager")
	TelegramId, err := strconv.ParseInt(pager, 10, 64)
	if err != nil {
		return err
	}
	u.TelegramId = TelegramId
	return nil
}
