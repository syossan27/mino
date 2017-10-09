package main

// Bufferはtermboxの内部バッファの責務を負う


// Size: 内部バッファの最大サイズ（コマンド履歴配列のサイズとも言い換えれる）
// Offset: 表示されるコマンド履歴の開始位置
// CommandPosition: コマンドの表示位置
type Buffer struct {
	Size            int
	Offset          int
	CommandPosition int
}

const (
	defaultCommandPosition = 1
	defaultBufferOffset    = 0
)

func NewBuffer(size int) Buffer {
	return Buffer {
		Offset:          defaultBufferOffset,
		CommandPosition: defaultCommandPosition,
		Size:            size,
	}
}

func (b *Buffer) ClearCommandPosition() {
	b.CommandPosition = defaultCommandPosition
}

func (b *Buffer) ClearOffset() {
	b.Offset = defaultBufferOffset
}
