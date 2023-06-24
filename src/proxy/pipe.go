package proxy

import (
	"io"
	"net"

	"golang.org/x/sync/errgroup"
)

type Pipe struct {
	Src net.Conn
	Dst net.Conn
}

func (p *Pipe) Close() {
	p.Src.Close()
	p.Dst.Close()
}

func NewPipe(src, dst net.Conn) *Pipe {
	return &Pipe{
		Src: src,
		Dst: dst,
	}
}

func (p *Pipe) Run() error {
	errGroup := errgroup.Group{}
	errGroup.Go(func() error {
		_, err := p.pipe(p.Dst, p.Src)
		return err
	})
	errGroup.Go(func() error {
		_, err := p.pipe(p.Src, p.Dst)
		return err
	})
	return errGroup.Wait()
}

func (p *Pipe) pipe(dst net.Conn, src net.Conn) (int64, error) {
	return io.Copy(dst, src)
}
