package radius

type Config struct {
	Debug         bool
	Address       string
	Port          uint16
	Secret        string
	Answertimeout int
}
