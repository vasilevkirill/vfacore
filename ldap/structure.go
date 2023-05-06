package ldap

type Config struct {
	User     string
	Password string
	Servers  []string
	Dn       string
}

type User struct {
	TelegramId     int64
	DispalyName    string
	SAMAccountName string
}
