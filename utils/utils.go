package utils

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
	"syscall"
)

const (
	maxDatagramSize = 8192
)

type HostStatus string

const (
	Query  string = "query"
	Lock   string = "lock"
	UnLock string = "unlock"

	Free   HostStatus = "free"
	Locked HostStatus = "locked"
)

type HostTwin struct {
	sync.Mutex
	Status HostStatus
}

// NewBroadcaster creates a new UDP multicast connection on which to broadcast
func NewBroadcaster() (*net.UDPConn, error) {
	broadcastAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:2000")
	if err != nil {
		panic(err)
	}
	udpConn, err := net.DialUDP("udp", nil, broadcastAddr)
	if err != nil {
		return nil, err
	}

	return udpConn, nil
}

// Listen binds to the UDP address and port given and writes packets received
// from that address to a buffer which is passed to a hander
func Listen(address string, handler func(*net.UDPConn, *net.UDPAddr, int, []byte)) {
	// Parse the string address
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		log.Fatal(err)
	}

	// Open up a connection
	conn, err := net.ListenUDP("udp4", addr)
	if err != nil {
		log.Fatal(err)
	}

	conn.SetReadBuffer(maxDatagramSize)
	defer conn.Close()
	// Loop forever reading from the socket
	for {
		buffer := make([]byte, maxDatagramSize)
		numBytes, src, err := conn.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}

		handler(conn, src, numBytes, buffer)
	}
}

func GracefullyExit() {
	fmt.Println("Start Exit...")
	fmt.Println("Execute Clean...")
	fmt.Println("End Exit...")
	os.Exit(0)
}

func RegisterSignal(c chan os.Signal) {
	go func() {
		for s := range c {
			switch s {
			case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
				fmt.Println("Program Exit...", s)
				GracefullyExit()
			case syscall.SIGUSR1:
				fmt.Println("usr1 signal", s)
			case syscall.SIGUSR2:
				fmt.Println("usr2 signal", s)
			default:
				fmt.Println("other signal", s)
			}
		}
	}()
}
