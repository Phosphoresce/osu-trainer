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
	var opt int

	// log file
	f, err := os.Create("ot.log")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	log.SetOutput(f)

	// init ncurses window
	text, err := gc.Init()
	if err != nil {
		log.Println("Init: ", err)
	}
	defer gc.End()

	// ncurses options
	gc.Cursor(0)
	gc.Raw(true)
	gc.Echo(false)

	// menu
	text.Clear()
	text.Keypad(true)
	menu_items := []string{"Test your skill", "140 bpm", "160 bpm", "180 bpm", "200 bpm", "Exit"}
	items := make([]*gc.MenuItem, len(menu_items))
	for i, val := range menu_items {
		items[i], _ = gc.NewItem(val, "")
		defer items[i].Free()
	}
	menu , err := gc.NewMenu(items)
	if err != nil {
		text.Print(err)
	}
	defer menu.Free()
	menu.Post()
	for {
		gc.Update()
		ch := text.GetChar()

		switch gc.KeyString(ch) {
		case "q":
			os.Exit(0)
		case "down":
			menu.Driver(gc.REQ_DOWN)
			opt++
		case "up":
			menu.Driver(gc.REQ_UP)
			opt--
		case "return", "enter":
			text.Erase()
			text.Keypad(false)
			switch opt {
			case 0:
				testYourSkill(text)
			case 1:
				testBpm(text, 140.00)
			case 2:
				testBpm(text, 160.00)
			case 3:
				testBpm(text, 180.00)
			case 4:
				testBpm(text, 200.00)
			case 5:
				os.Exit(0)
			}
		}
	}
}

func testYourSkill(text *gc.Window) {
	var zcount, xcount, row int
	var totaltaps, deviation, gap float64
	var grade string
	tpm := 0.0

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
		text.Refresh()

		handleInput(text, &zcount, &xcount, &row, &start, &lastTap, &gap)
		totaltaps = float64(zcount + xcount)
		tpm = totaltaps/time.Now().Sub(start).Minutes()
		deviation = 1000/(tpm/60)
		if gap < deviation+25 && gap > deviation-25 {
			grade = "good"
		} else {
			grade = "bad"
		}
		text.MovePrintf(7 + row, 0, "%8.2f ms [%s]", gap, grade)
		row++
		if row >= 50 {
			row = 0
		}
	}
}

func testBpm(text *gc.Window, bpm float64) {
	var zcount, xcount, row int
	var totaltaps, deviation, gap float64
	var grade string
	tpm := 0.0
	deviation = 1000/((bpm*4)/60)

	// start the clock
	start := time.Now()
	lastTap := start

	for {
		// updateConsole
		text.MovePrintf(0, 0, "tpm: %6.2f z: %6d x: %6d", tpm, zcount, xcount)
		text.MovePrintf(2, 0, "%v bpm mode", bpm)
		text.MovePrintf(3, 0, "Target between taps: %5.2f(+/-25) ms", deviation)
		text.MovePrintf(5, 0, "Current 1/4 note stream speed: %6.2f bpm", tpm/4)
		text.MovePrintf(7, 0, "Time between taps:")
		text.Refresh()

		handleInput(text, &zcount, &xcount, &row, &start, &lastTap, &gap)
		totaltaps = float64(zcount + xcount)
		tpm = totaltaps/time.Now().Sub(start).Minutes()
		if gap < deviation+25 && gap > deviation-25 {
			grade = "good"
		} else {
			grade = "bad"
		}
		text.MovePrintf(8 + row, 0, "%8.2f ms [%s]", gap, grade)
		row++
		if row >= 50 {
			row = 0
		}
	}
}

func handleInput(text *gc.Window, zcount, xcount, row *int, start, lastTap *time.Time, gap *float64) {
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
		*row = 0
		text.Erase()
	}
}
