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

type OperateCode string
type HostStatus string

const (
	Query  OperateCode = "query"
	Lock   OperateCode = "lock"
	UnLock OperateCode = "unlock"

	Free   HostStatus = "free"
	Locked HostStatus = "locked"
)

type UdpMessage struct {
	Code OperateCode
}

type HostTwin struct {
	sync.Mutex
	Status HostStatus
}

// NewBroadcaster creates a new UDP multicast connection on which to broadcast
func NewBroadcaster(port int) (*net.UDPConn, error) {
	lAddr := &net.UDPAddr{
		IP:   net.IPv4(127, 0, 0, 1),
		Port: port,
	}

	// 这里设置接收者的IP地址为广播地址
	rAddr := &net.UDPAddr{
		IP:   net.IPv4(255, 255, 255, 255),
		Port: port,
	}
	conn, err := net.DialUDP("udp", lAddr, rAddr)
	if err != nil {
		return nil, err
	}

	return conn, nil
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
