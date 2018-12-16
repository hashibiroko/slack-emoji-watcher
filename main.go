package main

import (
	"golang.org/x/net/websocket"
	"log"
	"net/http"
	"encoding/json"
	"fmt"
	"flag"
	"os"
)

const rtm = "https://slack.com/api/rtm.connect"
const OriginalUrl = "http://localhost:8787"

type Flags struct {
	SlackBotToken string
	ChannelName   string
}

type SlackRTMConnect struct {
	OK  bool   `json:"ok"`
	URL string `json:"url"`
}

type Receive struct {
	Type  string
	Error Error
}

type Error struct {
	Code int32
	Msg  string
}

func main() {
	flags := loadFlags()

	req, _ := http.NewRequest("GET", rtm, nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	query := req.URL.Query()
	query.Add("token", flags.SlackBotToken)
	req.URL.RawQuery = query.Encode()

	client := new(http.Client)
	res, _ := client.Do(req)
	defer res.Body.Close()

	var slackRTMConnect SlackRTMConnect
	decoder := json.NewDecoder(res.Body)
	err := decoder.Decode(&slackRTMConnect)
	if err != nil {
		log.Fatal(err)
	}

	if slackRTMConnect.OK {
		ws, err := websocket.Dial(slackRTMConnect.URL, "", OriginalUrl)
		if err != nil {
			log.Fatal(err)
		}
		var jsonRes interface{}
		for {
			websocket.JSON.Receive(ws, &jsonRes)

			var emojiChanged = jsonRes.(map[string]interface{})
			if emojiChanged["type"] == "emoji_changed" &&
				emojiChanged["subtype"] == "add" {
				sendMessage(flags, emojiChanged["name"].(string))
			}
		}
	}
}

func sendMessage(flag Flags, emojiName string) {
	req, _ := http.NewRequest("POST", "https://slack.com/api/chat.postMessage", nil)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	query := req.URL.Query()
	query.Add("token", flag.SlackBotToken)
	query.Add("channel", flag.ChannelName)
	query.Add("as_user", "true")
	query.Add("text", fmt.Sprintf("あなた、もしかして :%s: ？", emojiName))
	req.URL.RawQuery = query.Encode()

	client := new(http.Client)
	res, _ := client.Do(req)
	defer res.Body.Close()
}

func loadFlags() Flags {

	slackBotToken := flag.String("token",
		os.Getenv("SLACK_BOT_TOKEN"),
		"Set your slack bot token.")

	channelName := flag.String("channel",
		os.Getenv("SLACK_CHANNEL_NAME"),
		"Set your slack notification channel.")

	flag.Parse()

	if *slackBotToken == "" {
		log.Println("Slack token is require.")
		os.Exit(1)
	}

	if *channelName == "" {
		*channelName = "random"
		log.Println("Set notification channel to random.")
	}

	flags := Flags{
		SlackBotToken: *slackBotToken,
		ChannelName:   *channelName,
	}

	return flags
}
