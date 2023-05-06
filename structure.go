package vfacore

type tgConfig struct {
	Debug          bool
	Token          string
	PoolAddress    string
	PoolPort       uint16
	HookDomain     string
	HookPort       uint16
	HookCertPub    string
	HookCertKey    string
	NameBot        string `yaml:"-"`
	WebHookAddress string `yaml:"-"`
}

type ldapConfig struct {
	User     string
	Password string
	Servers  []string
	Dn       string
}

type ldapUser struct {
	TelegramId     int64
	DispalyName    string
	SAMAccountName string
}

type radiusConfig struct {
	Debug         bool
	Address       string
	Port          uint16
	Secret        string
	Answertimeout int
	ServerAddress string `yaml:"-"`
}

type configS struct {
	Telegram tgConfig     `yaml:"telegram"`
	Ldap     ldapConfig   `yaml:"ldap"`
	Radius   radiusConfig `yaml:"radius"`
}
