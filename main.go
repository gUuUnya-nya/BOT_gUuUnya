package main

import (
	"context"
	"fmt"
	"github.com/traPtitech/go-traq"
	"strings"

	"github.com/traPtitech/traq-ws-bot"
	"github.com/traPtitech/traq-ws-bot/payload"
)

func main() {
	bot, err := traqwsbot.NewBot(&traqwsbot.Options{
		AccessToken: "bekVJISDicjnsfpddvTWRmu6CkkO9jLsNLLi",
	})
	if err != nil {
		panic(err)
	}

	bot.OnError(func(message string) {
		fmt.Println("Received ERROR message: " + message)
	})

	bot.OnDirectMessageCreated(func(p *payload.DirectMessageCreated) {
		messageText := p.Message.Text
		fmt.Println("Message Received: " + messageText)

		if !strings.HasPrefix(messageText, "!DM") {
			fmt.Println("Not start with !DM")
		}

		splitMessageText := strings.SplitN(messageText, "$", 3)
		if len(splitMessageText) != 3 {
			fmt.Println("splitMessageText length is not 3")
		}

		userIds := splitMessageText[1]
		dm := splitMessageText[2]
		fmt.Println(dm)
		recipientList := strings.Split(userIds, ", ")
		fmt.Println(recipientList)

		for _, recipient := range recipientList {
			users, response, err := bot.API().UserApi.GetUsers(context.Background()).Name(recipient).Execute()
			if err != nil || response.StatusCode != 200 {
				fmt.Println("GetUsers Error:" + err.Error())
			}

			if len(users) == 0 {
				fmt.Println("Error: User not found")
				continue
			} else if len(users) >= 2 {
				fmt.Println("Error: Multiple users found")
				continue
			}

			user := users[0]

			_, response, err = bot.API().MessageApi.PostDirectMessage(context.Background(), user.Id).PostMessageRequest(traq.PostMessageRequest{Content: dm}).Execute()
			if err != nil || response.StatusCode != 201 {
				fmt.Println("PostDirectMessage Error:" + err.Error())
			}
		}
	})

	if err := bot.Start(); err != nil {
		panic(err)
	}
}