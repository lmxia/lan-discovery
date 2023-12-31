package main

import (
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/lmxia/lan-discovery/utils"
	log "github.com/sirupsen/logrus"
)

var (
	myStatus utils.HostTwin
)

func msgHandler(conn *net.UDPConn, src *net.UDPAddr, n int, content []byte) {
	var err error
	message := string(content[:n])
	log.Infof("来自%v,Discovered %v", src, string(content[:n]))
	// 广播的查询接口
	if message == utils.Query {
		myStatus.Lock()
		defer myStatus.Unlock()
		_, err = conn.WriteToUDP([]byte(myStatus.Status), src) // 发送数据
		if err != nil {
			log.Infof("Write to udp failed, err:  %s", err)
		}
	} else if message == utils.Lock {
		myStatus.Lock()
		defer myStatus.Unlock()
		myStatus.Status = utils.Locked
		_, err = conn.WriteToUDP([]byte("ok, i was locked just now."), src) // 发送数据
		if err != nil {
			log.Infof("Write to udp failed, err:  %s", err)
		}
	} else if message == utils.UnLock {
		myStatus.Lock()
		defer myStatus.Unlock()
		myStatus.Status = utils.Free
		_, err = conn.WriteToUDP([]byte("ok, i am free right now."), src) // 发送数据
		if err != nil {
			log.Infof("Write to udp failed, err:  %s", err)
		}
	} else {
		_, err = conn.WriteToUDP([]byte("i can't understand."), src) // 发送数据
		if err != nil {
			log.Infof("Write to udp failed, err:  %s", err)
		}
	}
}

func main() {
	c := make(chan os.Signal)
	myStatus.Status = utils.Free
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM,
		syscall.SIGQUIT, syscall.SIGUSR1, syscall.SIGUSR2)
	utils.RegisterSignal(c)

	utils.Listen(2000, msgHandler)
}
