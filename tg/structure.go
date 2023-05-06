package tg

type Config struct {
	Debug       bool
	Token       string
	PoolAddress string
	PoolPort    uint16
	HookDomain  string
	HookPort    uint16
	HookCertPub string
	HookCertKey string
}
