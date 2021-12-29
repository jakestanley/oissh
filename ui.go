package main

import (
	"log"
	"sync/atomic"
	"time"

	ui "github.com/gizak/termui/v3"
	"github.com/gizak/termui/v3/widgets"
)

var (
	messages *widgets.Paragraph
	system   *widgets.Paragraph
	users    *widgets.Paragraph
)

func initUi() {

	if err := ui.Init(); err != nil {
		log.Fatalf("failed to initialize termui: %v", err)
	}

	// TODO auto resize, etc
	system = widgets.NewParagraph()
	system.Text = ""
	system.Title = "System"
	system.SetRect(0, 0, 80, 4)

	messages = widgets.NewParagraph()
	messages.Text = ""
	messages.Title = "Messages"
	messages.SetRect(0, 5, 80, 9)

	users = widgets.NewParagraph()
	users.Text = ""
	users.Title = "Users"
	users.SetRect(0, 10, 80, 14)
}

func writeUsersString(message string) {

	users.Text = message
}

func writeSystemString(message string) {

	system.Text = message
}

func writeMessageString(message string) {

	messages.Text = message
}

func renderUI() {

	for {
		ui.Render(system, messages, users)
		time.Sleep(time.Second / 2)
	}
}

func inputUi() {
	defer atomic.StoreInt64(&kill, 1)
	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			break
		}
		time.Sleep(time.Second)
	}
}
