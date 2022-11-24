package main

import (
	"fmt"
	"os"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	usage = `
 countdown

 Usage:
    countdown 25s
    countdown 1m50s
    countdown 3h40m50s

 Options:
    -up : count up from zero

`
	tick = time.Second
)

var (
	timer          *time.Timer
	ticker         *time.Ticker
	queues         chan termbox.Event
	startDone      bool
	startX, startY int
)

func draw(d time.Duration) {
	w, h := termbox.Size()
	clear()

	str := format(d)
	text := toText(str)

	if !startDone {
		startDone = true
		startX, startY = w/2-text.width()/2, h/2-text.height()/2
	}

	x, y := startX, startY
	for _, s := range text {
		echo(s, x, y)
		x += s.width()
	}

	flush()
}

func format(d time.Duration) string {
	d = d.Round(time.Second)
	h := d / time.Hour
	d -= h * time.Hour
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	if h < 1 {
		return fmt.Sprintf("%02d:%02d", m, s)
	}
	return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}

func start(d time.Duration) {
	timer = time.NewTimer(d)
	ticker = time.NewTicker(tick)
}

func stop() {
	timer.Stop()
	ticker.Stop()
}

func countdown(timeLeft time.Duration, countUp bool) {
	var exitCode int

	start(timeLeft)

	if countUp {
		timeLeft = 0
	}

	draw(timeLeft)

loop:
	for {
		select {
		case ev := <-queues:
			if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC) {
				exitCode = 1
				break loop
			}
			if ev.Ch == 'p' || ev.Ch == 'P' {
				stop()
			}
			if ev.Ch == 'c' || ev.Ch == 'C' {
				start(timeLeft)
			}
		case <-ticker.C:
			if countUp {
				timeLeft += tick
			} else {
				timeLeft -= tick
			}
			draw(timeLeft)
		case <-timer.C:
			break loop
		}
	}

	termbox.Close()
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

func getKitchenTimeDuration(date string) (time.Duration, error) {

	targetTime, err := time.Parse(time.Kitchen, date)
	if err != nil {
		return time.Duration(0), err
	}

	now := time.Now()
	originTime := time.Date(0, time.January, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

	// the time of day has already passed, so target tomorrow
	if targetTime.Before(originTime) {
		targetTime = targetTime.AddDate(0, 0, 1)
	}

	duration := targetTime.Sub(originTime)

	return duration, err
}

func main() {
	if len(os.Args) < 2 || len(os.Args) > 3 {
		stderr(usage)
		os.Exit(2)
	}

	term := os.Getenv("TERM")
	if (len(term) >= 4 && term[0:4] == "tmux") {
		os.Setenv("TERM", "xterm-color")
	}

	timeLeft, err := getKitchenTimeDuration(os.Args[1])

	if err != nil {

		timeLeft, err = time.ParseDuration(os.Args[1])
		if err != nil {
			stderr("error: invalid duration or kitchen time: %v\n", os.Args[1])
			os.Exit(2)
		}
	}

	err = termbox.Init()
	if err != nil {
		panic(err)
	}

	queues = make(chan termbox.Event)
	go func() {
		for {
			queues <- termbox.PollEvent()
		}
	}()
	countUp := len(os.Args) == 3 && os.Args[2] == "-up"
	countdown(timeLeft, countUp)
}
