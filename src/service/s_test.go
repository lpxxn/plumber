package service

import (
	"bufio"
	"fmt"
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

type myIClienter interface {
}
type myTClient struct {
	Name string
}

func (c *myTClient) Run(i myIClienter) error {
	switch i.(type) {
	case *myTClient:
		t := i.(*myTClient)
		t.Name = "hello"
		return nil
	default:
		if i == nil {
			fmt.Println("interface is nil")
			return nil
		}
		fmt.Println("interface is not nil")

		return nil
	}
}

func TestInterfaceFunc(t *testing.T) {
	m := &myTClient{}
	m.Run(nil)
	t.Log(m.Name)

	m.Run("")
	t.Log()
}

func TestBf(t *testing.T) {
	// bufio write data to os.Stdout
	// 默认的太大了，会导致数据不及时输出
	writer := bufio.NewWriter(os.Stdout)
	//writer := bufio.NewWriterSize(os.Stdout, 5)
	writer.WriteString("hello world!")
	writer.WriteString("abcdefg")
	//writer.Flush()
	exitChan := make(chan struct{})

	select {
	case <-exitChan:
		// Channel is not closed
	default:
		// Channel is closed
	}
	close(exitChan)

	//	close(exitChan) // error
	if _, ok := <-exitChan; !ok {
		t.Log("Channel is closed")
	} else {
		t.Log("Channel is not closed")
	}

	exitChan2 := make(chan struct{}, 1)
	close(exitChan2)
	// close(exitChan2)
}
