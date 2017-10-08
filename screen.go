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

var (
	color = map[string]termbox.Attribute {
		"default": termbox.ColorDefault,
		"white": termbox.ColorWhite,
		"black": termbox.ColorBlack,
		"magenta": termbox.ColorMagenta,
		"cyan": termbox.ColorCyan,
	}

	attr = map[string]termbox.Attribute {
		"underline": termbox.AttrUnderline,
	}

	style = map[string]termbox.Attribute {
		"fg": color["white"],
		"bg": color["default"],
		"selectionFg": color["white"] | attr["underline"],
		"selectionBg": color["magenta"],
		"selectedFg": color["black"],
		"selectedBg": color["cyan"],
	}
)

func NewTermbox() *Termbox {
	s := NewSelection()
	b := NewBuffer()

	return &Termbox{
		Selection: s,
		Buffer: b,
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
	buffer := t.Buffer
	selection := t.Selection

	termbox.Clear(color["default"], color["default"])

	// Debug
	tso := strconv.Itoa(t.Selection.Offset)
	tbo := strconv.Itoa(t.Buffer.Offset)
	th := strconv.Itoa(t.Height)
	lens := strconv.Itoa(len(selection.Selected))
	t.Print(0, 0, color["white"], color["default"], "Selection Offset: " + tso + ", Buffer Offset: " + tbo + ", Terminal Height: " + th + ", len(s.Selected): " + lens)

	// コマンド履歴表示開始位置から順にコマンド履歴を表示する
	// indexはコマンド履歴配列の添字, buffer.Offsetを初期値とする
	// commandYはコマンドの表示場所
	// selection.Offsetはコマンド履歴を選択している表示位置
	// commandIndexはコマンド全体
	for i := buffer.Offset; i < len(commandHistory); i++ {
		command := commandHistory[i]
		// これの名前変えたい
		commandIndex := buffer.commandY + buffer.Offset

		// 選択済みかどうか
		selectedNumber := ""
		number, exist := selection.GetSelectedIndex(commandIndex)
		if exist {
			selectedNumber = strconv.Itoa(number) + ": "

			// 要リファクタリング
			if commandIndex == selection.Offset {
				t.PrintSelection(0, buffer.commandY, style["selectionFg"], style["selectionBg"], selectedNumber + command)
			} else {
				t.PrintSelection(0, buffer.commandY, style["selectedFg"], style["selectedBg"], selectedNumber + command)
			}
		} else {
			if commandIndex == selection.Offset {
				t.PrintSelection(0, buffer.commandY, style["selectionFg"], style["selectionBg"], command)
			} else {
				t.Print(0, buffer.commandY, style["fg"], style["bg"], command)
			}
		}

		buffer.commandY++
	}
	termbox.Flush()
	buffer.ClearCommandY()
}

func (t *Termbox) Print(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func (t *Termbox) PrintSelection(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}

	// 横幅いっぱいに選択状態の背景色を表示する
	for ; x < t.Width; x++ {
		termbox.SetCell(x, y, ' ', fg, bg)
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
			case termbox.KeySpace:
				t.Selection.Select(commandHistory)
			case termbox.KeyCtrlD:
				t.Selection.ClearSelected()
			}
		}

		t.Draw(commandHistory)
	}
}
