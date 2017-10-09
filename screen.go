package main

import (
	"strconv"
	"github.com/nsf/termbox-go"
	"github.com/mattn/go-runewidth"
)

type (
	Termbox struct {
		Width int
		Height int
		SearchFormHeight int
		CommandHistoryHeight int
		CommandHistory []Command
		Selection
		Buffer
		Filter
	}
)

var (
	color = map[string]termbox.Attribute {
		"default": termbox.ColorDefault,
		"white": termbox.ColorWhite,
		"black": termbox.ColorBlack,
		"magenta": termbox.ColorMagenta,
		"cyan": termbox.ColorCyan,
		"red": termbox.ColorRed,
		"green": termbox.ColorGreen,
		"blue": termbox.ColorBlue,
	}

	attr = map[string]termbox.Attribute {
		"underline": termbox.AttrUnderline,
	}

	style = map[string]termbox.Attribute {
		"fg": color["white"],
		"bg": color["default"],
		"searchFg": color["cyan"],
		"searchBg": color["black"],
		"selectionFg": color["white"] | attr["underline"],
		"selectionBg": color["magenta"],
		"selectedFg": color["black"],
		"selectedBg": color["green"],
		"selectedSearchFg": color["magenta"],
		"selectedSearchBg": color["green"],
	}
)

func NewTermbox(commandHistory []Command) *Termbox {
	return &Termbox{
		Selection: NewSelection(),
		Buffer: NewBuffer(len(commandHistory)),
		Filter: NewFilter(),
		CommandHistory: commandHistory,
	}
}

func (t *Termbox) Init() error {
	err := termbox.Init()
	if err != nil {
		return err
	}

	return nil
}

func (t *Termbox) Display() {
	defer termbox.Close()

	termbox.SetInputMode(termbox.InputAlt | termbox.InputMouse)

	t.SetSize()
	t.Draw()

loop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyArrowUp:
				// TODO: もうちょっとスッキリかける気がする・・・？

				// カーソルが検索フォームに被らないように制御
				if t.Selection.Index > 1 {
					t.Selection.Index--
				}

				// カーソルを最上段から一つ上に動かした時に、
				// コマンド履歴の開始位置を下げる
				if t.Selection.Index < t.SearchFormHeight + t.Buffer.Offset {
					t.Buffer.Offset--
				}
			case termbox.KeyArrowDown:
				// TODO: もうちょっとスッキリかける気がする・・・？

				// カーソルがコマンド履歴分を越えないよう制御
				if t.Selection.Index < t.Buffer.Size {
					t.Selection.Index++
				}

				// カーソルを最下段から一つ下に動かした時に、
				// コマンド履歴の開始位置を上げる
				if t.Selection.Index > t.CommandHistoryHeight + t.Buffer.Offset {
					t.Buffer.Offset++
				}
			case termbox.KeyCtrlC:
				break loop
			case termbox.KeyEsc:
				break loop
			case termbox.KeySpace:
				t.Selection.Select()
			case termbox.KeyCtrlD:
				t.Selection.ClearSelected()
			case termbox.KeyBackspace2:
				if t.Filter.Length() > 0 {
					t.Filter.Delete()
				}

				// 位置関係を初期化
				t.Selection.ClearIndex()
				t.Buffer.ClearOffset()
			case termbox.KeyEnter:
			default:
				if ev.Ch != 0 {
					t.Filter.Append(ev.Ch)

					// 位置関係を初期化
					t.Selection.ClearIndex()
					t.Buffer.ClearOffset()
				}
			}
		}

		t.Draw()
	}
}

func (t *Termbox) SetSize() {
	w, h := termbox.Size()
	t.Width = w
	t.Height = h
	t.SearchFormHeight = 1
	t.CommandHistoryHeight = h - t.SearchFormHeight
}

func (t *Termbox) Draw() {
	termbox.Clear(color["default"], color["default"])

	commandHistory := t.CommandHistory

	// 検索フォームの表示
	displaySearchForm := "QUERY>"
	t.Print(0, 0, style["fg"], style["bg"], displaySearchForm)
	t.PrintSearchQuery(len(displaySearchForm), 0, style["fg"], style["bg"])

	// 検索条件がある場合、検索条件に合致するコマンド履歴一覧を生成する
	// コマンド履歴配列が変更されるため、内部バッファのサイズも更新
	if len(t.Filter.SearchQuery) != 0 {
		commandHistory = t.Filter.FilterResult(t.CommandHistory)
		t.Buffer.Size = len(commandHistory)
	}

	// commandHistoryを順番に表示
	for i := t.Buffer.Offset; i < len(commandHistory); i++ {
		command := commandHistory[i]

		// 選択済みかどうか
		number, exist := t.Selection.GetSelectedNumber(command.Index)

		// 表示させる文字列
		var appendExecOrder string
		var displayCommand string
		var displayColorSet map[string]termbox.Attribute
		var displaySearchColorSet map[string]termbox.Attribute
		if exist {
			// 選択済の場合
			appendExecOrder = strconv.Itoa(number) + ": "
			displayCommand = appendExecOrder + command.Content
			displayColorSet = map[string]termbox.Attribute {
				"fg": style["selectedFg"],
				"bg": style["selectedBg"],
			}
			displaySearchColorSet = map[string]termbox.Attribute {
				"fg": style["selectedSearchFg"],
				"bg": style["selectedSearchBg"],
			}
		} else {
			// 未選択の場合
			displayCommand = command.Content
			displayColorSet = map[string]termbox.Attribute {
				"fg": style["fg"],
				"bg": style["bg"],
			}
			displaySearchColorSet = map[string]termbox.Attribute {
				"fg": style["searchFg"],
				"bg": style["searchBg"],
			}
		}

		// コマンド履歴の表示
		if t.Buffer.CommandPosition == t.Selection.Index - t.Buffer.Offset {
			// カーソルの表示
			t.PrintFullWidth(0, t.Buffer.CommandPosition, style["selectionFg"], style["selectionBg"], displayCommand)

			// 検索条件に合致した箇所に着色
			// カーソル表示時は着色する色は固定
			t.Print(len(appendExecOrder) + command.FilterIndex, t.Buffer.CommandPosition, color["cyan"], style["selectionBg"], string(t.Filter.SearchQuery))

			// カーソルにあるコマンドを仮保存
			t.Selection.Command = command
		} else {
			// 選択済か、未選択のコマンド履歴の表示
			t.PrintFullWidth(0, t.Buffer.CommandPosition, displayColorSet["fg"], displayColorSet["bg"], displayCommand)

			// 検索条件に合致した箇所に着色
			t.Print(len(appendExecOrder) + command.FilterIndex, t.Buffer.CommandPosition, displaySearchColorSet["fg"], displaySearchColorSet["bg"], string(t.Filter.SearchQuery))
		}

		t.Buffer.CommandPosition++
	}

	termbox.Flush()
	t.Buffer.ClearCommandPosition()
}

func (t *Termbox) Print(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		width := runewidth.RuneWidth(c)
		x = x + width
	}
}

func (t *Termbox) PrintFullWidth(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		width := runewidth.RuneWidth(c)
		x = x + width
	}

	// 横幅いっぱいに背景色を表示する
	for ; x < t.Width; x++ {
		termbox.SetCell(x, y, ' ', fg, bg)
	}
}

func (t *Termbox) PrintSearchQuery(x, y int, fg, bg termbox.Attribute) {
	for _, c := range t.Filter.SearchQuery {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}
