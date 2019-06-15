package marvin

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"
	"time"

	"github.com/awgh/goirc/logging"
	"github.com/awgh/markov"
	irc "github.com/fluffle/goirc/client"
)

// Run - Main run method - infinite loop
func Run(configDir string, chainFile string, mcflyFile string) {

	logging.SetLogger(sLogger{})

	InitMarkovChains(chainFile, mcflyFile)

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

			// if SlackAPIToken is defined, we assume it's a Slack connection and ignore everything else
			if len(ircClientConfig.SlackAPIToken) > 0 {
				startSlackClient(&ircClientConfig, db)
			} else {
				startIrcClient(&ircClientConfig, db)
			}
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
