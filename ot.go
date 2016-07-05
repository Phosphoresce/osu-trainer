package main

import (
	"os"
	"log"
	"time"
	gc "github.com/rthornton128/goncurses"
)

func main() {
	_main()
}

func _main() {
	var zcount, xcount int
	var totaltaps, deviation, gap float64
	var grade string
	tpm := 0.0

	// init ncurses menu here
	f, err := os.Create("ot.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	log.SetOutput(f)

	text, err := gc.Init()
	if err != nil {
		log.Println("Init: ", err)
	}
	defer gc.End()

	// ncurses options
	gc.Cursor(0)
	gc.Raw(true)
	gc.Echo(false)

	// start the clock
	start := time.Now()
	lastTap := start

	for {
		// updateConsole
		text.MovePrintf(0, 0, "tpm: %6.2f z: %6d x: %6d", tpm, zcount, xcount)
		text.MovePrintf(2, 0, "You can handle:")
		text.MovePrintf(3, 0, "1/2 note streams at %6.2f bpm", tpm/2)
		text.MovePrintf(4, 0, "1/4 note streams at %6.2f bpm", tpm/4)
		text.MovePrintf(6, 0, "Time between taps:")
		text.MovePrintf(7 + zcount + xcount, 0, "%8.2f ms [%s]", gap, grade)
		text.Refresh()

		handleInput(text, &zcount, &xcount, &start, &lastTap, &gap)
		totaltaps = float64(zcount + xcount)
		tpm = totaltaps/time.Now().Sub(start).Minutes()
		deviation = 1000/(tpm/60)
		if gap < deviation+25 && gap > deviation-25 {
			grade = "good"
		} else {
			grade = "bad"
		}
	}
}

func handleInput(text *gc.Window, zcount, xcount *int, start, lastTap *time.Time, gap *float64) {
	key := text.GetChar()

	switch byte(key) {
	case 'q':
		os.Exit(0)
	case 'z':
		*zcount++
		*gap = time.Since(*lastTap).Seconds()*1000
		*lastTap = time.Now()
	case 'x':
		*xcount++
		*gap = time.Since(*lastTap).Seconds()*1000
		*lastTap = time.Now()
	case 'r':
		*start = time.Now()
		*lastTap = *start
		*gap = 0
		*zcount = 0
		*xcount = 0
		text.Erase()
	}
}
