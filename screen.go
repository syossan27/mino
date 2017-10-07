package main

import (
	"fmt"
	"github.com/nsf/termbox-go"
)

const coldef = termbox.ColorDefault

var (
	current string
	curev termbox.Event
)

type Termbox struct {}

func NewTermbox() *Termbox {
	return &Termbox{}
}

func (t *Termbox) Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}

	return nil
}

func (t *Termbox) Draw() {
	termbox.Clear(coldef, coldef)
	t.Tbprint(0, 0, termbox.ColorMagenta, coldef, "test hoge fuga")
	termbox.Flush()
}

func (t *Termbox) Tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func (t *Termbox) Do() {
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)
	t.Draw()

	data := make([]byte, 0, 64)
mainloop:
	for {
		if cap(data)-len(data) < 32 {
			newdata := make([]byte, len(data), len(data)+32)
			copy(newdata, data)
			data = newdata
		}
		beg := len(data)
		d := data[beg : beg+32]
		switch ev := termbox.PollRawEvent(d); ev.Type {
		case termbox.EventRaw:
			data = data[:beg+ev.N]
			current = fmt.Sprintf("%q", data)
			if current == `"q"` {
				break mainloop
			}

			for {
				ev := termbox.ParseEvent(data)
				if ev.N == 0 {
					break
				}
				curev = ev
				copy(data, data[curev.N:])
				data = data[:len(data)-curev.N]
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		t.Draw()
	}
}
