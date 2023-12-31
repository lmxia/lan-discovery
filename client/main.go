package main

import (
	"fmt"
	"net"
	"os"
	"os/signal"

	"github.com/lmxia/lan-discovery/utils"
	log "github.com/sirupsen/logrus"
)

// UDP 客户端
func main() {
	c := make(chan os.Signal)
	signal.Notify(c)

	socket, dstAddr, err := utils.NewBroadcaster()
	if err != nil {
		log.Infof("UDP广播发送失败，%s: ", err)
	}
	defer socket.Close()

	go listenUdp(socket)

	_, err = socket.WriteToUDP([]byte(utils.Query), dstAddr) // 发送数据
	if err != nil {
		fmt.Println("发送数据失败，err: ", err)
		return
	}
	<-c
}

func listenUdp(socket *net.UDPConn) {
	for {
		data := make([]byte, 4096)
		n, remoteAddr, err := socket.ReadFromUDP(data) // 接收数据
		log.Infof("read %s from <%s>\n", data[:n], remoteAddr)
		if err != nil {
			fmt.Println("接收数据失败, err: ", err)
			return
		}

		ack := string(data[:n])
		if ack == string(utils.Free) {
			_, err = socket.WriteToUDP([]byte(utils.Lock), remoteAddr) // 发送数据
			if err != nil {
				fmt.Println("发送lock数据失败，err: ", err)
				return
			}
		} else if ack == string(utils.Locked) {
			_, err = socket.WriteToUDP([]byte(utils.UnLock), remoteAddr) // 发送数据
			if err != nil {
				fmt.Println("发送unlock数据失败，err: ", err)
				return
			}
		}
	}
}
