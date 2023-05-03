package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	usage = `
 countdown [-up] [-say] <duration>

 Usage
  countdown 25s
  countdown 14:15
  countdown 02:15PM

 Flags
`
	tick         = time.Second
	inputDelayMS = 500 * time.Millisecond
)

var (
	timer          *time.Timer
	ticker         *time.Ticker
	queues         chan termbox.Event
	w, h           int
	inputStartTime time.Time
	isPaused       bool
)

func main() {
	countUp := flag.Bool("up", false, "count up from zero")
	sayTheTime := flag.Bool("say", false, "announce the time left")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		stderr(usage)
		flag.PrintDefaults()
		os.Exit(2)
	}
	timeLeft, err := parseTime(args[0])

	if err != nil {
		timeLeft, err = time.ParseDuration(args[0])
		if err != nil {
			stderr("error: invalid duration or time: %v\n", args[0])
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
	countdown(timeLeft, *countUp, *sayTheTime)
}

func start(d time.Duration) {
	timer = time.NewTimer(d)
	ticker = time.NewTicker(tick)
}

func stop() {
	timer.Stop()
	ticker.Stop()
}

func durationToDraw(timeLeft, totalDuration time.Duration, countUp bool) time.Duration {
	if countUp {
		return totalDuration - timeLeft
	}
	return timeLeft
}

func countdown(totalDuration time.Duration, countUp bool, sayTheTime bool) {
	timeLeft := totalDuration
	var exitCode int
	isPaused = false
	w, h = termbox.Size()
	start(timeLeft)

	draw(durationToDraw(timeLeft, totalDuration, countUp), w, h)
	if sayTheTime {
		go say(timeLeft)
	}

loop:
	for {
		select {
		case ev := <-queues:
			if ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC {
				exitCode = 1
				break loop
			}

			if pressTime := time.Now(); ev.Key == termbox.KeySpace && pressTime.Sub(inputStartTime) > inputDelayMS {
				if isPaused {
					start(timeLeft)
					draw(durationToDraw(timeLeft, totalDuration, countUp), w, h)
				} else {
					stop()
					drawPause(w, h)
				}

				isPaused = !isPaused
				inputStartTime = time.Now()
			}

			if ev.Type == termbox.EventResize {
				w, h = termbox.Size()
				draw(durationToDraw(timeLeft, totalDuration, countUp), w, h)

				if isPaused {
					drawPause(w, h)
				}
			}
		case <-ticker.C:
			timeLeft -= tick
			draw(durationToDraw(timeLeft, totalDuration, countUp), w, h)
			if sayTheTime {
				go say(timeLeft)
			}
		case <-timer.C:
			break loop
		}
	}

	termbox.Close()
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}

func draw(d time.Duration, w int, h int) {
	clear()

	str := format(d)
	text := toText(str)

	startX, startY := w/2-text.width()/2, h/2-text.height()/2

	x, y := startX, startY
	for _, s := range text {
		echo(s, x, y)
		x += s.width()
	}

	flush()
}

func drawPause(w int, h int) {
	startX := w/2 - pausedText.width()/2
	startY := h * 3 / 4

	echo(pausedText, startX, startY)
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

func say(d time.Duration) {
	if d.Seconds() <= 10 {
		cmd := exec.Command("say", fmt.Sprintf("%v", d.Seconds()))
		_ = cmd.Run()
	}
}

func parseTime(date string) (time.Duration, error) {
	targetTime, err := time.Parse(time.Kitchen, strings.ToUpper(date))
	if err != nil {
		targetTime, err = time.Parse("15:04", date)
		if err != nil {
			return time.Duration(0), err
		}
	}

	now := time.Now()
	originTime := time.Date(0, time.January, 1, now.Hour(), now.Minute(), now.Second(), 0, time.UTC)

	// The time of day has already passed, so target tomorrow.
	if targetTime.Before(originTime) {
		targetTime = targetTime.AddDate(0, 0, 1)
	}

	duration := targetTime.Sub(originTime)

	return duration, err
}
