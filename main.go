package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/caarlos0/env"
	"github.com/slack-go/slack"
)

type config struct {
	SlackToken   string `env:"SLACK_TOKEN"`
	GetChannelId string `env:"GET_CHANNEL_ID"`
	Since        string `env:"SINCE"`
	Until        string `env:"UNTIL"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Panic(err)
	}

	today, lw := getLastWeek()

	var since time.Time
	if cfg.Since == "" {
		since = lw
	} else {
		var err error
		since, err = time.Parse(time.RFC3339, cfg.Since)
		if err != nil {
			log.Panic(err)
		}
	}

	var until time.Time
	if cfg.Until == "" {
		until = today
	} else {
		var err error
		until, err = time.Parse(time.RFC3339, cfg.Until)
		if err != nil {
			log.Panic(err)
		}
	}

	from, to := strconv.FormatInt(since.Unix(), 10), strconv.FormatInt(until.Unix(), 10)

	res, err := getConversations(cfg.SlackToken, cfg.GetChannelId, from, to)
	if err != nil {
		log.Panic(err)
	}

	fmt.Println(res)
}

var nowFunc func() time.Time

func nowUTC() time.Time {
	if nowFunc != nil {
		return nowFunc()
	}
	return time.Now().UTC()
}

func getLastWeek() (time.Time, time.Time) {
	y, m, d := nowUTC().Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, time.UTC)

	lw := today.AddDate(0, 0, -7)

	return today, lw
}

func getConversations(slackToken, channelId, from, to string) ([]slack.Message, error) {
	api := slack.New(slackToken)

	params := slack.GetConversationHistoryParameters{ChannelID: channelId, Oldest: from, Latest: to, Limit: 1000}
	conv, err := api.GetConversationHistory(&params)
	if err != nil {
		return nil, err
	}

	if !conv.SlackResponse.Ok {
		return nil, fmt.Errorf("error response: %s", conv.SlackResponse.Error)
	}

	return conv.Messages, nil
}
