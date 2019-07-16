package main

import (
	"context"
	"sync"
	"flag"
	"github.com/hiank/think/net"
	"github.com/hiank/thinkend/master"
	"github.com/golang/glog"
	// "github.com/hiank/think/db"
)

func main() {

	defer glog.Infoln("close master-server")

	flag.Parse()

	wg := new(sync.WaitGroup)
	wg.Add(2)

	net.Init(context.Background())
	go serveK8s(wg)
	go serveWS(wg)

	wg.Wait()
}

func serveK8s(wg *sync.WaitGroup) {

	defer wg.Done()

	if err := net.ServeK8s("", &master.Handler{}); err != nil {

		glog.Fatalln("serve k8s error : " + err.Error())
	}
}


func serveWS(wg *sync.WaitGroup) {

	defer wg.Done()

	if err := net.ServeWS(""); err != nil {

		glog.Fatalln("serve ws error : " + err.Error())
	}
}