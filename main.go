package main

import (
	"fmt"
	"log"
	"syscall"
	"time"
)

func main() {
	fd, err := syscall.Socket(syscall.AF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Print(fd)

	sa := &syscall.SockaddrInet4{
		Port: 9999,
		Addr: [4]byte{0, 0, 0, 0},
	}

	if err := syscall.Bind(fd, sa); err != nil {
		log.Fatal(err)
	}

	if err := syscall.Listen(fd, 0); err != nil {
		log.Fatal(err)
	}

	nfd, _, err := syscall.Accept(fd)
	if err != nil {
		log.Fatal(err)
	}

	b := make([]byte, 1024)
 	n, err := syscall.Read(nfd, b)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("bytes read: %d", n)
	
	time.Sleep(time.Hour)
}
