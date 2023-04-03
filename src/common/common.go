package common

import "sync"

const (
	MagicString = "Go"
)

var SeparatorBytes = []byte(" ")

type WaitGroup struct {
	sync.WaitGroup
}

func (w *WaitGroup) WaitFunc(f func()) {
	w.Add(1)
	go func() {
		f()
		w.Done()
	}()
}
