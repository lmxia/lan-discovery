package utils

import (
	"fmt"
	"net"
	"os"
	"sync"
	"syscall"

	log "github.com/sirupsen/logrus"
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
func NewBroadcaster() (*net.UDPConn, *net.UDPAddr, error) {
	srcAddr := &net.UDPAddr{IP: net.IPv4zero, Port: 0}
	dstAddr := &net.UDPAddr{IP: net.ParseIP("192.168.0.255"), Port: 2000}
	conn, err := net.ListenUDP("udp", srcAddr)
	if err != nil {
		return nil, nil, err
	}

	return conn, dstAddr, nil
}

// Listen binds to the UDP address and port given and writes packets received
// from that address to a buffer which is passed to a hander
func Listen(port int, handler func(*net.UDPConn, *net.UDPAddr, int, []byte)) {
	listener, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.IPv4zero, Port: port})
	if err != nil {
		fmt.Println(err)
		return
	}
	log.Infof("Local: <%s> \n", listener.LocalAddr().String())
	listener.SetReadBuffer(maxDatagramSize)
	defer listener.Close()
	// Loop forever reading from the socket
	buffer := make([]byte, maxDatagramSize)
	for {
		numBytes, src, err := listener.ReadFromUDP(buffer)
		if err != nil {
			log.Fatal("ReadFromUDP failed:", err)
		}

		handler(listener, src, numBytes, buffer)
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
