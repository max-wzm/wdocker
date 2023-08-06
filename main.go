package main

import (
	"os"
	"time"
	"wdocker/log"

	"github.com/max-wzm/geerpc"
)

const (
	usage = `mydocker is a simple container runtime implementation.
	The purpose of this project is to learn how docker works and how to write a docker by ourselves
	Enjoy it, just for fun.`
)

func main() {
	raddr := geerpc.StartRegistry(9999)
	f, err := os.OpenFile("registry", os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		log.Error(err.Error())
		return
	}
	f.WriteString(raddr)
	var daemon Daemon
	geerpc.StartServer(raddr, &daemon)
	time.Sleep(time.Second * 100)
}
