package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/caarlos0/env"
)

type config struct {
	Since string `env:"SINCE"`
	Until string `env:"UNTIL"`
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
	fmt.Println(since, until)
	fmt.Println(from, to)
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
