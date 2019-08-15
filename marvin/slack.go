package marvin

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/nlopes/slack"
)

func startSlackClient(config *MarvinConfig, db *sql.DB) {

	// You more than likely want your "Bot User OAuth Access Token" which starts with "xoxb-"
	api := slack.New(
		config.SlackAPIToken,
		slack.OptionDebug(true),
		slack.OptionLog(log.New(os.Stdout, "slack-bot: ", log.Lshortfile|log.LstdFlags)),
	)

	rtm := api.NewRTM()
	go rtm.ManageConnection()

	go func() {
		for msg := range rtm.IncomingEvents {
			fmt.Print("Event Received: ")
			switch ev := msg.Data.(type) {
			case *slack.HelloEvent:
				// Ignore hello

			case *slack.ConnectedEvent:
				fmt.Println("Infos:", ev.Info)
				fmt.Println("Connection counter:", ev.ConnectionCount)
				/*
					channelID, timestamp, err := api.PostMessage(config.SlackChannel,
						slack.MsgOptionText("I think you ought to know, I'm feeling rather depressed.", false)) //slack.MsgOptionAttachments(attachment))
					if err != nil {
						fmt.Printf("%s\n", err)
						return
					}
					fmt.Printf("Message successfully sent to channel %s at %s", channelID, timestamp)
				*/
			case *slack.MessageEvent:
				fmt.Printf("Message: %v\n", ev)

				handleSlack(api, ev, config, db)

			case *slack.PresenceChangeEvent:
				fmt.Printf("Presence Change: %v\n", ev)

			case *slack.LatencyReport:
				fmt.Printf("Current latency: %v\n", ev.Value)

			case *slack.RTMError:
				fmt.Printf("Error: %s\n", ev.Error())

			case *slack.InvalidAuthEvent:
				fmt.Printf("Invalid credentials")
				return

			default:
				// Ignore other events..
				// fmt.Printf("Unexpected: %v\n", msg.Data)
			}
		}
	}()
}
