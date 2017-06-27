package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"database/sql"

	irc "github.com/fluffle/goirc/client"
	"github.com/fluffle/goirc/logging"

	_ "github.com/mattn/go-sqlite3"
)

type sLogger struct{}

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

// Message contains an answering machine message
type Message struct {
	From   string
	To     string
	Text   string
	Public bool
}

// LinePrinter prints a line for debugging
func LinePrinter(line *irc.Line) {
	log.Println("Public:", line.Public())
	log.Println("Target:", line.Target())
	log.Println("Text:", line.Text())
	log.Println("Args:", line.Args)
	log.Println("Cmd:", line.Cmd)
	log.Println("Host:", line.Host)
	log.Println("Ident:", line.Ident)
	log.Println("Nick:", line.Nick)
	log.Println("Raw:", line.Raw)
	log.Println("Src:", line.Src)
	log.Println("Tags:", line.Tags)
	log.Println("Time:", line.Time)
}

func remove(s []string, r string) []string {
	for i, v := range s {
		if v == r {
			return append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func (s sLogger) Debug(f string, a ...interface{}) {
	log.Printf(f, a...)
}
func (s sLogger) Info(f string, a ...interface{}) {
	log.Printf(f, a...)
}
func (s sLogger) Warn(f string, a ...interface{}) {
	log.Printf(f, a...)
}
func (s sLogger) Error(f string, a ...interface{}) {
	log.Printf(f, a...)
}

var configDir = flag.String("confdir", "config", "Config Directory")
var ircClients []*irc.Conn

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

	c.HandleFunc(irc.PRIVMSG,
		func(conn *irc.Conn, line *irc.Line) {

			LinePrinter(line)

			deliverMessages(conn, line.Nick, config.Channel)

			var sendFn func(string)
			if line.Public() { // respond to public messages publicly
				sendFn = func(msg string) { conn.Notice(config.Channel, msg) }
			} else { // respond to private messages privately
				sendFn = func(msg string) { conn.Privmsg(line.Nick, msg) }
			}
			broadcastFn := func(msg string) { conn.Privmsg(config.Channel, msg) }
			sendPriv := func(msg string) { conn.Notice(line.Nick, msg) }

			args := strings.Split(line.Args[1], " ")
			if len(args) > 0 {
				switch args[0] {

				case ".h":
					fallthrough
				case ".help":
					sendPriv("*****" + string(2) + "Marvin Help" + string(0xF) + "*****")
					sendPriv(string(2) + "Marvin " + string(0xF) + "responds to private messages privately and responds to channel commands as notices,")
					sendPriv("with the exception of the .5questions command, where the response is always broadcast to the channel.")
					sendPriv("The following commands are available:")
					sendPriv(string(2) + ".5 [username]" + string(0xF) + "(alias: .5questions)")
					sendPriv(" will broadcast the Five Questions, with an optional greeting for " +
						string(2) + "username" + string(0xF) + " to the channel.")
					sendPriv(string(2) + ".booze [booze_name_or_prefix]" + string(0xF) + "(alias: .b)")
					sendPriv(" will list Boozes used in the mixed drink database.  This works as a string prefix search.")
					sendPriv("If there is more than one match, all matches will be listed.  If no argumet is given, all Boozes will be listed.")
					sendPriv("If only one Booze matches, the list of Drinks using that Booze will be shown.")
					sendPriv(string(2) + ".drink [drink_name_or_prefix]" + string(0xF) + "(alias: .d)")
					sendPriv(" will display Drink recipes from the mixed drink database.  This works as a string prefix search.")
					sendPriv("If there is more than one match, all matches will be listed.  If no argumet is given, all Drinks will be listed.")
					sendPriv("If only one Drink matches, the recipe for that drink will be shown.")
					sendPriv(string(2) + ".tell <nick> <message>" + string(0xF) + "(alias: .t)")
					sendPriv(" will send a message to nick the next time they join or talk in channel.  Private tells will be sent privately.")
					break

				case ".5":
					fallthrough
				case ".5questions":
					if len(args) > 1 {
						broadcastFn("Greetings " + string(2) + args[1] + string(0xF) + " and Welcome to " + string(2) + "Milliways" + string(0xF) + ", the Restaurant at the End of the Universe!")
					}
					broadcastFn("  Please answer the following questions, by way of introduction:")
					broadcastFn("  1.  Who are you?")
					broadcastFn("  2.  How did you get here?")
					broadcastFn("  3.  What can Milliways do for you?")
					broadcastFn("  4.  What can you do for Milliways?")
					broadcastFn("  5.  What are you good at that isn't computers?")
					break

				case ".macker":
					sendFn("macker is a twat")
					break

				case ".d":
					fallthrough
				case ".drink":
					if len(args) > 1 {
						drinkName := strings.Join(args[1:], "%") + "%"
						log.Println("DRINKNAME: " + drinkName)
						rc, err := db.Query("SELECT name, ingredients, prep from drinks WHERE name LIKE ? COLLATE NOCASE;", drinkName)
						if err != nil {
							log.Fatal(err.Error())
						}
						var dnames []string
						var dname, ingredients, prep string
						for rc.Next() {
							rc.Scan(&dname, &ingredients, &prep)
							dnames = append(dnames, dname)
						}
						if len(dnames) == 1 {
							sendFn(string(2) + dname)
							ilines := strings.Split(ingredients, "\n")
							for _, il := range ilines {
								sendFn(il)
							}
							sendFn(prep)
						} else if len(dnames) > 0 {
							sendFn(strings.Join(dnames, ", "))
						}
					} else {
						rc, err := db.Query("SELECT DISTINCT name from drinks ORDER BY name ASC;")
						if err != nil {
							log.Fatal(err.Error())
						}
						var dks []string
						for rc.Next() {
							var drink string
							rc.Scan(&drink)
							dks = append(dks, drink)
						}
						if len(dks) > 0 {
							sendFn(strings.Join(dks, ", "))
						}
					}
					break

				case ".b":
					fallthrough
				case ".booze":
					if len(args) > 1 {
						boozeName := strings.Join(args[1:], "%") + "%"
						log.Println("BOOZENAME: " + boozeName)

						rc, err := db.Query("SELECT ingredient, drink from ingredients WHERE ingredient LIKE ? COLLATE NOCASE;", boozeName)
						if err != nil {
							log.Fatal(err.Error())
						}
						boozes := make(map[string]int)
						var dks []string
						var bname, drink string
						for rc.Next() {
							rc.Scan(&bname, &drink)
							dks = append(dks, drink)
							if _, ok := boozes[bname]; !ok {
								boozes[bname] = 1
							}
						}
						keys := []string{}
						for k := range boozes {
							keys = append(keys, k)
						}
						if len(keys) == 1 && len(dks) > 0 {
							sendFn("Drinks made with " + string(2) + bname + string(0xF) + ": " + strings.Join(dks, ", "))
						} else if len(keys) > 1 {
							sendFn(strings.Join(keys, ", "))
						}
					} else {
						rc, err := db.Query("SELECT DISTINCT ingredient from ingredients ORDER BY ingredient ASC;")
						if err != nil {
							log.Fatal(err.Error())
						}
						var igs []string
						for rc.Next() {
							var ingredient string
							rc.Scan(&ingredient)
							igs = append(igs, ingredient)
						}
						if len(igs) > 0 {
							sendFn(strings.Join(igs, ", "))
						}
					}
					break

				case ".t":
					fallthrough
				case ".tell":
					if len(args) > 3 {
						_, ok := namesMessages[args[1]]
						if !ok {
							namesMessages[line.Nick] = make([]Message, 1)
						}
						namesMessages[args[1]] = append(namesMessages[args[1]], Message{
							From:   line.Nick,
							To:     args[1],
							Text:   strings.Join(args[2:], " "),
							Public: line.Public(),
						})
						sendFn("fine, I will relay your message... here I am, brain the size of a planet...")
					}
					break

				default: // The Wormhole Case : forward public messages across servers
					for i := range ircClients {
						if ircClients[i] != conn {
							ircClients[i].Privmsg(config.Channel, line.Nick+" "+strings.Join(args[:], " "))
						}
					}
					break
				}
			}

		})

	c.HandleFunc(irc.JOIN,
		func(conn *irc.Conn, line *irc.Line) {
			log.Println("JOIN: " + line.Nick)
			hostsChansNames[config.Host][config.Channel] = append(hostsChansNames[config.Host][config.Channel], line.Nick)
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
			//log.Println(line)
		})

	/*
	   352    RPL_WHOREPLY
	                 "<channel> <user> <host> <server> <nick>
	                 ( "H" / "G" > ["*"] [ ( "@" / "+" ) ]
	                 :<hopcount> <real name>"

	*/
	c.HandleFunc("352",
		func(conn *irc.Conn, line *irc.Line) {
			log.Println("WHO:", line.Args[5])
			hostsChansNames[config.Host][config.Channel] = append(hostsChansNames[config.Host][config.Channel], line.Args[5])
		})

	ircClients = append(ircClients, c)
	return nil
}

// Map of hostnames to a map of chan names to a list of string nicks
var hostsChansNames map[string]map[string][]string

// Map of answering machine messages
var namesMessages map[string][]Message

func main() {
	flag.Parse()
	logging.SetLogger(sLogger{})

	hostsChansNames = make(map[string]map[string][]string)
	namesMessages = make(map[string][]Message)

	// Drinks DB setup
	db, err := sql.Open("sqlite3", "./IBA-Cocktails-2016.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	files, err := ioutil.ReadDir(*configDir) // open all json files in this directory, parse them, and call startIrcClient with them
	if err != nil {
		log.Fatal(err)
	}
	for _, file := range files {
		fmt.Println(file.Name())

		if strings.HasSuffix(file.Name(), ".json") {
			dat, err := ioutil.ReadFile(*configDir + string(os.PathSeparator) + file.Name())
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
