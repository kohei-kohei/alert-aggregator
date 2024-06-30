package main

import (
	"fmt"
	"log"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/caarlos0/env"
	"github.com/slack-go/slack"
)

type config struct {
	SlackToken    string `env:"SLACK_TOKEN,required"`
	GetChannelId  string `env:"GET_CHANNEL_ID,required"`
	SendChannelId string `env:"SEND_CHANNEL_ID"`
	Since         string `env:"SINCE"`
	Until         string `env:"UNTIL"`
}

func main() {
	cfg := config{}
	if err := env.Parse(&cfg); err != nil {
		log.Panic(err)
	}

	since, until, err := getAggregationPeriod(cfg.Since, cfg.Until)
	if err != nil {
		log.Panic(err)
	}

	from, to := strconv.FormatInt(since.Unix(), 10), strconv.FormatInt(until.Unix(), 10)

	res, err := getConversations(cfg.SlackToken, cfg.GetChannelId, from, to)
	if err != nil {
		log.Panic(err)
	}

	alerts := aggregateAlerts(res)

	keys := make([]string, 0, len(alerts))
	for k := range alerts {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return alerts[keys[i]] > alerts[keys[j]]
	})

	var alertContent strings.Builder
	var total int
	for _, k := range keys {
		v := alerts[k]
		fmt.Fprintf(&alertContent, "%s: %d\n", k, v)
		total += v
	}

	var output strings.Builder
	fmt.Fprintf(&output, "%s 〜 %s\n", since, until)
	fmt.Fprintf(&output, "total number of alerts: %d\n", total)
	fmt.Fprintf(&output, "number of alert types: %d\n\n", len(alerts))
	output.WriteString(alertContent.String())

	fmt.Print(output.String())

	if cfg.SendChannelId != "" {
		if err := sendToSlack(cfg.SlackToken, cfg.SendChannelId, output.String()); err != nil {
			log.Panic(err)
		}
	}
}

var nowFunc func() time.Time

func now() time.Time {
	if nowFunc != nil {
		return nowFunc()
	}
	return time.Now()
}

func getLastWeek() (time.Time, time.Time) {
	y, m, d := now().Date()
	today := time.Date(y, m, d, 0, 0, 0, 0, time.Local)

	lw := today.AddDate(0, 0, -7)

	return today, lw
}

func getAggregationPeriod(sinceStr, untilStr string) (time.Time, time.Time, error) {
	today, lw := getLastWeek()

	var since time.Time
	if sinceStr == "" {
		since = lw
	} else {
		var err error
		since, err = time.Parse(time.RFC3339, sinceStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	var until time.Time
	if untilStr == "" {
		until = today
	} else {
		var err error
		until, err = time.Parse(time.RFC3339, untilStr)
		if err != nil {
			return time.Time{}, time.Time{}, err
		}
	}

	if since.After(until) {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid period: 'since' (%s) is after 'until' (%s)", since, until)
	}

	return since, until, nil
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

func aggregateAlerts(messages []slack.Message) map[string]int {
	m := map[string]int{}

	for _, msg := range messages {
		if msg.BotProfile == nil {
			continue
		}

		bn := msg.BotProfile.Name
		if bn == "incoming-webhook" {
			continue
		}

		var alertTitle string
		if bn == "Datadog" {
			alertTitle = msg.Attachments[0].Title

			if strings.Contains(alertTitle, "Recovered") {
				continue
			}

			sep := ":"
			if strings.Contains(alertTitle, sep) {
				alertTitle = strings.Split(alertTitle, sep)[1]
			}
		} else if bn == "digdag-alert" {
			alertTitle = msg.Attachments[0].Title
		} else if bn == "AWS Chatbot" {
			alertTitle = msg.Attachments[0].Fallback

			if !strings.Contains(alertTitle, ":rotating_light:") {
				continue
			}

			sep := " | "
			if strings.Contains(alertTitle, sep) {
				alertTitle = strings.Split(alertTitle, sep)[1]
			}
		} else if bn == "Sentry" {
			alertTitle = msg.Blocks.BlockSet[0].(*slack.SectionBlock).Text.Text

			sep := "*"
			if strings.Contains(alertTitle, sep) {
				alertTitle = strings.Split(alertTitle, sep)[1]
			}
		} else if bn == "PagerDuty" {
			alertTitle = msg.Blocks.BlockSet[0].(*slack.SectionBlock).Text.Text

			if !strings.Contains(alertTitle, ":large_green_circle:") {
				continue
			}

			reg := regexp.MustCompile(`\|(.*?)>\*`)
			match := reg.FindStringSubmatch(alertTitle)

			if len(match) > 1 {
				alertTitle = match[1]
			}
		} else {
			log.Printf("This message from %s is not supported\n", bn)
			continue
		}

		key := bn + " | " + alertTitle
		m[key] = m[key] + 1
	}

	return m
}

func sendToSlack(slackToken, channelId, text string) error {
	api := slack.New(slackToken)

	_, _, err := api.PostMessage(channelId, slack.MsgOptionText(text, false))
	return err
}
