package main

// Bufferはtermboxの内部バッファの責務を負う
// Offsetは「表示されるコマンド履歴の開始位置」を担う

type Buffer struct {
	searchFormY int
	commandY    int
	Offset      int
}

const (
	initSearchFormY  = 0
	initCommandY     = 1
	initBufferOffset = 0
)

func NewBuffer() Buffer {
	return Buffer {
		searchFormY: initSearchFormY,
		commandY:    initCommandY,
		Offset:      initBufferOffset,
	}
}

func (b *Buffer) ClearCommandY() {
	b.commandY = initCommandY
}
