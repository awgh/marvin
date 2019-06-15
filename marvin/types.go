package marvin

// Map of answering machine messages
var namesMessages map[string][]Message

// Map of hostnames to a map of chan names to a list of string nicks
var hostsChansNames map[string]map[string][]string

func init() {
	namesMessages = make(map[string][]Message)
	hostsChansNames = make(map[string]map[string][]string)
}

// MarvinConfig - Client Config for Marvin
type MarvinConfig struct {
	Host         string
	Port         string
	Nick         string
	Password     string
	Name         string
	Version      string
	Quit         string
	ProxyEnabled bool
	SASL         bool
	Proxy        string
	Channel      string

	MD5ApiUser string
	MD5ApiCode string

	SlackAPIToken string
	SlackChannel  string
}

// Message - Contains an answering machine message
type Message struct {
	From   string
	To     string
	Text   string
	Public bool
}
