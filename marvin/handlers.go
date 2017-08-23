package main

import (
	"database/sql"
	"log"
	"strings"

	irc "github.com/fluffle/goirc/client"

	_ "github.com/mattn/go-sqlite3"
)

func handlePrivMsg(conn *irc.Conn, line *irc.Line, config *MarvinConfig, db *sql.DB) {

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
			sendPriv(string(2) + ".5questions [username]" + string(0xF) + "(alias: .5)")
			sendPriv(" will broadcast the Five Questions, with an optional greeting for " +
				string(2) + "username" + string(0xF) + " to the channel.")
			sendPriv(string(2) + ".x4questions [username]" + string(0xF) + "(alias: .x4)")
			sendPriv(" will ask additional four Questions.")
			sendPriv(string(2) + ".booze [booze_name_or_prefix]" + string(0xF) + "(alias: .b)")
			sendPriv(" will list Boozes used in the mixed drink database.  This works as a string prefix search.")
			sendPriv("If there is more than one match, all matches will be listed.  If no argumet is given, all Boozes will be listed.")
			sendPriv("If only one Booze matches, the list of Drinks using that Booze will be shown.")
			sendPriv(string(2) + ".drink [drink_name_or_prefix]" + string(0xF) + "(alias: .d)")
			sendPriv(" will display Drink recipes from the mixed drink database.  This works as a string prefix search.")
			sendPriv("If there is more than one match, all matches will be listed.  If no argument is given, all Drinks will be listed.")
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

		case ".x4":
			fallthrough
		case ".x4questions":
			if len(args) > 1 {
				broadcastFn("Hi " + string(2) + args[1] + string(0xF))
			}
			broadcastFn("  Here are some extra questions, by lubiana and macker:")
			broadcastFn("  1.  Are you cute?")
			broadcastFn("  2.  Do you like cuddles?")
			broadcastFn("  3.  Do you like cute dogs or cats?")
			broadcastFn("  4.  How much time a day do you spend cuddling?")
			break

		case ".macker":
			sendFn("macker is a twat")
			break

		case ".mcfly":
			// maybe switch to chain file
			// http://www.imdb.com/character/ch0001829/quotes
			sendFn("If you put your mind to it, you can accomplish anything.")
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
			if len(args) > 2 {
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

		case ".m":
			if loadedChain != nil {
				sendFn(loadedChain.Generate(23))
			} else {
				log.Println("No chain file loaded")
			}

		default: // The Wormhole Case : forward public messages across servers
			for i := range ircClients {
				if ircClients[i] != conn {
					ircClients[i].Privmsg(config.Channel, "@_"+line.Nick+" "+strings.Join(args[:], " "))
				}
			}
			break
		}
	}
}
