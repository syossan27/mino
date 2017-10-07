package main

import (
	"strconv"
	"github.com/nsf/termbox-go"
)

type Termbox struct {
	Width int
	Height int
	Selection
	Buffer
}

type Selection struct {
	Offset int
}

type Buffer struct {
	Offset int
}

func NewTermbox() *Termbox {
	return &Termbox{
		Selection: Selection {
			Offset: 1,
		},
		Buffer: Buffer {
			Offset: 0,
		},
	}
}

func (t *Termbox) Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}

	return nil
}

func (t *Termbox) SetSize() {
	w, h := termbox.Size()
	t.Width = w
	t.Height = h - 1
}

func (t *Termbox) Draw(commandHistory []string) {
	y := 1
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	tso := strconv.Itoa(t.Selection.Offset)
	tbo := strconv.Itoa(t.Buffer.Offset)
	th := strconv.Itoa(t.Height)
	t.Tbprint(0, 0, termbox.ColorWhite, termbox.ColorDefault, "Selection Offset: " + tso + ", Buffer Offset: " + tbo + ", Terminal Height: " + th)

	for i := t.Buffer.Offset; i < len(commandHistory); i++ {
		command := commandHistory[i]
		if y + t.Buffer.Offset == t.Selection.Offset {
			t.Tbprint(0, y, termbox.ColorWhite, termbox.ColorMagenta, command)
		} else {
			t.Tbprint(0, y, termbox.ColorWhite, termbox.ColorDefault, command)
		}
		y++
	}
	termbox.Flush()
}

func (t *Termbox) Tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func (t *Termbox) Do(commandHistory []string) {
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)

	t.SetSize()
	t.Draw(commandHistory)

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				if t.Selection.Offset > 1 {
					t.Selection.Offset = t.Selection.Offset - 1
				}

				if t.Selection.Offset < t.Buffer.Offset + 1 {
					t.Buffer.Offset--
				}
			case termbox.KeyArrowDown:
				t.Selection.Offset = t.Selection.Offset + 1
				if t.Selection.Offset > t.Height + t.Buffer.Offset {
					t.Buffer.Offset++
				}
			case termbox.KeyCtrlC:
				break mainloop
			case termbox.KeyEsc:
				break mainloop
			}
		}

		t.Draw(commandHistory)
	}
}
