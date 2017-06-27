package main

import (
	"log"

	"github.com/awgh/markov"
	irc "github.com/fluffle/goirc/client"
)

func deliverMessages(conn *irc.Conn, nick string, channel string) {
	v, ok := namesMessages[nick]
	if ok {
		for msg := range v {
			if v[msg].Public {
				conn.Notice(channel, v[msg].From+" left a message for "+v[msg].To+": "+v[msg].Text)
			} else {
				conn.Privmsg(v[msg].To, v[msg].From+" left a message: "+v[msg].Text)
			}
		}
		namesMessages[nick] = make([]Message, 0)
	}
}

func loadChainFile(filename string) (*markov.Chain, bool) {
	chain := markov.NewChain(2) //prefix length should likely always be 2
	if err := chain.Load(filename); err != nil {
		log.Println(err)
		return nil, false
	}
	log.Println("Loaded chain file: " + filename)
	return chain, true
}

func seenNick(nick string, config *MarvinConfig) {
	hostsChansNames[config.Host][config.Channel] = append(hostsChansNames[config.Host][config.Channel], nick)
}

func remove(s []string, r string) []string { // todo: move to global utils
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}
