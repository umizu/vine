package vine

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Vine struct {
	errCh     chan error // stops server on receive
	stopDelay int        // gracefully shutdown, in milliseconds
}

func New() *Vine {
	return &Vine{
		errCh:     make(chan error, 1),
		stopDelay: 1500,
	}
}

func (v *Vine) Start() error { // todo: add param 'listenAddr'
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return err
	}

	if err := syscall.SetsockoptInt(fd, syscall.SOL_SOCKET, syscall.SO_REUSEADDR, 1); err != nil {
		return err
	}

	sa := &syscall.SockaddrInet4{
		Port: 9999,
		Addr: [4]byte{0, 0, 0, 0},
	}

	if err := syscall.Bind(fd, sa); err != nil {
		return err
	}

	if err := syscall.Listen(fd, 0); err != nil {
		return err
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go v.acceptLoop(fd)

	select {
	case <-sigs:
		syscall.Close(fd)
		time.Sleep(time.Millisecond * time.Duration(v.stopDelay))
		break
	case e := <-v.errCh:
		err = e
		break
	}

	return err
}

func (v *Vine) acceptLoop(sockfd int) {
	for {
		nfd, _, err := syscall.Accept(sockfd)
		if err != nil {
			slog.Debug("accept error", "err", err.Error())
			continue
		}

		go v.handleConn(nfd)
	}
}

func (v *Vine) handleConn(fd int) {
	req := parseRequest(fd)

	_ = req // todo: handle

	if err := syscall.Close(fd); err != nil {
		slog.Error("close err", "err", err.Error())
	}
}
