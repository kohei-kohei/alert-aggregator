package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println(getLastWeek())
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
