package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/pansapiens/go-defrag-sim/scandisk"

	"github.com/gdamore/tcell/v2"
	"github.com/spf13/pflag"
)

func main() {
	var skipScan bool

	// Parse command-line arguments
	pflag.IntVarP(&scandisk.Width, "width", "w", 0, "Width of the terminal")
	pflag.IntVarP(&scandisk.Height, "height", "h", 0, "Height of the terminal")
	pflag.IntVarP(&scandisk.BaseDelay, "delay", "d", 100, "Base delay in milliseconds")
	pflag.BoolVarP(&scandisk.Loop, "loop", "l", false, "Run in an infinite loop")
	pflag.BoolVar(&skipScan, "skip-scan", false, "Skip the scanning step")

	pflag.Parse()

	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err := screen.Init(); err != nil {
		panic(err)
	}
	defer screen.Fini()

	// Set Width and Height based on command-line arguments or terminal size
	screenWidth, screenHeight := screen.Size()
	if scandisk.Width == 0 {
		scandisk.Width = screenWidth
	}
	if scandisk.Height == 0 {
		scandisk.Height = screenHeight - 2 // Adjust for the status and legend bars
	}

	screen.Clear()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		screen.Fini()
		os.Exit(0)
	}()

	go func() {
		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlC || ev.Rune() == 'q' {
					screen.Fini()
					os.Exit(0)
				}
			}
		}
	}()

	for {
		scandisk.RunDefrag(skipScan, screen)
		if !scandisk.Loop {
			break
		}
	}
	screen.Clear()
	screen.Fini()
}
