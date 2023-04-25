package tg

type Config struct {
	Debug       bool
	Token       string
	PoolAddress string
	PoolPort    uint16
	HookUrl     string
	HookPort    uint16
	HookCertPub string
	HookCertKey string
}
