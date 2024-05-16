package traqwsbot_test

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/traPtitech/go-traq"

	"github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

func Example() {
	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: os.Getenv("lAvvrxwlgfA4I90i5xNFb8sZ3FYrwmHh0Zml"), // Required
		Origin:      "wss://q.trap.jp",         // Optional (default: wss://q.trap.jp)
	})
	if err != nil {
		panic(err)
	}

	bot.OnError(func(message string) {
		log.Println("Received ERROR message: " + message)
	})
	bot.OnMessageCreated(func(p *payload.MessageCreated) {
		log.Println("Received MESSAGE_CREATED event: " + p.Message.Text)
		_, _, err := bot.API().
			MessageApi.
			PostMessage(context.Background(), p.Message.ChannelID).
			PostMessageRequest(traq.PostMessageRequest{
				Content: "oisu-",
			}).
			Execute()
		if err != nil {
			log.Println(err)
		}
	})
	bot.OnDirectMessageCreated(func(p *payload.DirectMessageCreated) {
		message := p.Message.Text
		if !strings.HasPrefix(message, "!DM"){
			return
		}
		parts := strings.SplitN(message, "#", 2)
		if len(parts) != 2 {
			return
		}
		userIds := strings.TrimSpace(parts[0][4:])
		dm := parts[1]
		recipientList := strings.Split(userIds, ", ")
		for _, recipient := range recipientList {
			channel, response, err := bot.API().UserApi.GetUserDMChannel(context.Background(), recipient).Execute()
			if err != nil || response.StatusCode != 200 {
				return
			}
			_, _, err = bot.API().MessageApi.PostMessage(context.Background(), channel.Id).PostMessageRequest(traq.PostMessageRequest{Content: dm}).Execute()
			if err != nil{
				log.Println(err)
			}
		}
		log.Println(p.Message.Text)
	})

	if err := bot.Start(); err != nil {
		panic(err)
	}
}
