package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	"database/sql"

	"github.com/awgh/markov"
	irc "github.com/fluffle/goirc/client"
	"github.com/fluffle/goirc/logging"

	_ "github.com/mattn/go-sqlite3"
)

func init() {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator for markov.
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
}

// Message - Contains an answering machine message
type Message struct {
	From   string
	To     string
	Text   string
	Public bool
}

var ircClients []*irc.Conn
var markovChains []*markov.Chain

// Map of hostnames to a map of chan names to a list of string nicks
var hostsChansNames map[string]map[string][]string

// Map of answering machine messages
var namesMessages map[string][]Message

func main() {

	var configDir, chainFile, mcflyFile string
	flag.StringVar(&configDir, "confdir", "config", "Config Directory")
	flag.StringVar(&chainFile, "chain", "markov.chain", "Markov Chain File")
	flag.StringVar(&mcflyFile, "mcfly", "mcfly.chain", "McFly Chain File")
	flag.Parse()
	logging.SetLogger(sLogger{})

	hostsChansNames = make(map[string]map[string][]string)
	namesMessages = make(map[string][]Message)

	// Markov Chain setup
	markovChains = make([]*markov.Chain, 2)
	markovChains[0], _ = loadChainFile(chainFile)
	markovChains[1], _ = loadChainFile(mcflyFile)

	// Drinks DB setup
	db, err := sql.Open("sqlite3", "./IBA-Cocktails-2016.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	files, err := ioutil.ReadDir(configDir) // open all json files in this directory, parse them, and call startIrcClient
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())

		if strings.HasSuffix(file.Name(), ".json") {
			dat, err := ioutil.ReadFile(configDir + string(os.PathSeparator) + file.Name())
			if err != nil {
				log.Fatal(err.Error())
			}
			var ircClientConfig MarvinConfig
			if err := json.Unmarshal(dat, &ircClientConfig); err != nil {
				log.Fatal(err.Error())
			}
			startIrcClient(&ircClientConfig, db)
		}
	}

	for {
		for _, v := range ircClients {
			if !v.Connected() {
				log.Println("Connecting...")
				if err := v.Connect(); err != nil {
					log.Println(err.Error())
				}
			}
		}
		time.Sleep(time.Minute / 4)
	}
}

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
