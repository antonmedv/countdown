package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
	"os"
	"time"
)

var (
	duration       time.Duration
	deadline       time.Time
	startDone      = false
	startX, startY int
)

func draw() {
	w, h := termbox.Size()
	clear()

	left := time.Until(deadline)
	str := format(left)
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

	if duration.Hours() < 1 {
		return fmt.Sprintf("%02d:%02d", m, s)
	} else {
		return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
	}
}

func main() {
	var err error
	var exitCode = 0

	if len(os.Args) != 2 {
		stderr(`usage: 
  countdown 25s
  countdown 1m50s
  countdown 2h45m50s
`)
		os.Exit(2)
	}

	duration, err = time.ParseDuration(os.Args[1])
	if err != nil {
		stderr("error: invalid duration: %v\n", os.Args[1])
		os.Exit(2)
	}

	deadline = time.Now().Add(duration)
	timeout := time.After(duration)

	err = termbox.Init()
	if err != nil {
		panic(err)
	}

	queues := make(chan termbox.Event)
	go func() {
		for {
			queues <- termbox.PollEvent()
		}
	}()

	draw()

loop:
	for {
		select {
		case ev := <-queues:
			if ev.Type == termbox.EventKey && (ev.Key == termbox.KeyEsc || ev.Key == termbox.KeyCtrlC) {
				exitCode = 1
				break loop
			}
		case <-timeout:
			break loop
		default:
			draw()
			time.Sleep(10 * time.Millisecond)
		}
	}

	termbox.Close()
	if exitCode != 0 {
		os.Exit(exitCode)
	}
}
