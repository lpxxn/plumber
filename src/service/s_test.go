package service

import (
	"bufio"
	"os"
	"testing"
	"time"
)

func TestTic(t *testing.T) {
	timeout1 := 250 * time.Millisecond
	testTicker := time.NewTicker(timeout1)
	defer testTicker.Stop()
	exitChan := make(chan struct{})
	go func() {
		time.Sleep(10 * time.Second)
		close(exitChan)
	}()
	t.Log("start", time.Now().Second())
	for {
		select {
		case <-testTicker.C:
			t.Log("testTicker", time.Now().Second())
		case <-time.After(timeout1):
			t.Log("after timeout1", time.Now().Second())
		case <-exitChan:
			t.Log("exitChan", time.Now().Second())
			goto exit
		}
	}
exit:
	t.Log("end")
}

func TestBf(t *testing.T) {
	// bufio write data to os.Stdout
	// 默认的太大了，会导致数据不及时输出
	writer := bufio.NewWriter(os.Stdout)
	//writer := bufio.NewWriterSize(os.Stdout, 5)
	writer.WriteString("hello world!")
	writer.WriteString("abcdefg")
	//writer.Flush()
}
