package main

import (
	"github.com/golang/glog"
	"github.com/hiank/think/net"
)

func main() {

	glog.Infoln(net.ServeK8s("", NewHandler(net.GetRuntime().Context)))
}
