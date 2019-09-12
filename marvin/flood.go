package marvin

import (
	"log"
	"strings"
	"time"

	irc "github.com/fluffle/goirc/client"
)

var (
	xEvents  int
	ySeconds int
)

type entry struct {
	count     int       // total hits since timestamp
	timestamp time.Time // timestamp of first hit
}

var nickFloodTable map[string]*entry

func init() {
	nickFloodTable = make(map[string]*entry)

	xEvents = 3
	ySeconds = 30
}

func checkNickFlood(conn *irc.Conn, config *MarvinConfig, nick string, ident string) {

	idx := strings.Index(ident, "!")
	if idx < 1 {
		log.Println("Ident misformed?", ident)
		return
	}
	banline := "*" + ident[idx:]

	log.Println("Checking Nick Flood on ", banline)

	c := conn.StateTracker().GetChannel(config.Channel)
	if c == nil {
		log.Println("Failed to get state tracker for channel")
		return
	}
	cp, ok := c.Nicks[nick]
	if !ok {
		log.Println("Failed to find nick: " + nick)

	} else {
		// exempt operators and admins and such from banhammer
		if cp.Admin || cp.Op || cp.Owner || cp.HalfOp {
			log.Println("User ", nick, " is exempt")
			return
		}
	}

	// garbage collection, anything older than Y seconds
	now := time.Now()
	for k, v := range nickFloodTable {
		if now.Sub(v.timestamp) > (time.Duration(ySeconds) * time.Second) {
			delete(nickFloodTable, k)
		}
	}
	if v, ok := nickFloodTable[banline]; ok {
		// entry exists, increment it
		v.count++
		if v.count > xEvents {
			log.Println("kickban time!")
			conn.Mode(config.Channel, "+b", banline)
			conn.Kick(config.Channel, nick, "My name is Nick. Nick Flood, Private Detective.")
		}

	} else {
		// entry does not exist, add it
		nickFloodTable[banline] = &entry{count: 1, timestamp: now}
	}

}
