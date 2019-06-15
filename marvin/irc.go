package marvin

import (
	"crypto/tls"
	"database/sql"
	"log"
	"time"

	irc "github.com/fluffle/goirc/client"
)

var ircClients []*irc.Conn

func startIrcClient(config *MarvinConfig, db *sql.DB) error {

	// create new IRC connection
	cfg := irc.NewConfig(config.Nick) //nick
	cfg.SSL = true
	cfg.SSLConfig = &tls.Config{ServerName: config.Host, InsecureSkipVerify: true} //todo: add in CA manually to avoid this
	//cfg.SASL = config.SASL
	cfg.Server = config.Host + ":" + config.Port
	//cfg.NewNick = func(n string) string { return n + "^" } // don't need this if register nick
	cfg.Me.Ident = config.Nick
	cfg.Me.Name = config.Name
	cfg.Pass = config.Password
	cfg.Version = config.Version
	cfg.QuitMessage = config.Quit

	if config.ProxyEnabled {
		cfg.Proxy = config.Proxy
	}
	cfg.PingFreq = time.Second * 120

	c := irc.Client(cfg)
	c.EnableStateTracking()

	c.HandleFunc(irc.CONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			hostsChansNames[config.Host] = make(map[string][]string) // initialize this host entry in the nick map
			conn.Join(config.Channel)                                //todo: vectorize this
		})
	c.HandleFunc(irc.DISCONNECTED,
		func(conn *irc.Conn, line *irc.Line) {
			log.Println("DISCONNECTED")
		})

	c.HandleFunc(irc.PRIVMSG, func(conn *irc.Conn, line *irc.Line) {
		handlePrivMsg(conn, line, config, db)
	})

	c.HandleFunc(irc.JOIN,
		func(conn *irc.Conn, line *irc.Line) {
			log.Println("JOIN: " + line.Nick)
			seenNick(line.Nick, config)
			deliverMessages(conn, line.Nick, config.Channel)
		})

	c.HandleFunc(irc.PART,
		func(conn *irc.Conn, line *irc.Line) {
			log.Println("PART: " + line.Nick)
			remove(hostsChansNames[config.Host][config.Channel], line.Nick)
		})

	c.HandleFunc(irc.PING,
		func(conn *irc.Conn, line *irc.Line) {
			//log.Println("PING:")
		})

	c.HandleFunc("352", // Bug: RPL_WHOREPLY does not populate line correctly, use Args[5] for nick
		func(conn *irc.Conn, line *irc.Line) {
			log.Println("WHO:", line.Args[5])
			seenNick(line.Args[5], config)
		})

	ircClients = append(ircClients, c)
	return nil
}
