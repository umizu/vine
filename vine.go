package vine

import (
	"fmt"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Vine struct {
	Listener net.Listener

	shutdownDelay int // graceful shutdown (after closing listener)
}

func New() *Vine {
	return &Vine{
		shutdownDelay: 1000,
	}
}

func (v *Vine) Start(address string) error {
	l, err := net.Listen("tcp", address)
	if err != nil {
		return err
	}
	fmt.Printf("vine listening on %s\n", address)

	v.Listener = l
	go v.acceptLoop()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs

	if err := v.Listener.Close(); err != nil {
		return err
	}

	time.Sleep(time.Millisecond * time.Duration(v.shutdownDelay))
	return nil
}

func (v *Vine) acceptLoop() {
	for {
		conn, err := v.Listener.Accept()
		if err != nil {
			slog.Debug("accept error", "err", err.Error())
			continue
		}

		go v.handleConn(conn)
	}
}

func (v *Vine) handleConn(conn net.Conn) {
	defer conn.Close()
	req := parseRequest(conn)
	_ = req // todo: handle
}
