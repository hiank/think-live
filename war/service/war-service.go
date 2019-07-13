package main

import (
	"context"
	"sync"
	"flag"
	"github.com/golang/glog"
	// "fmt"
	"github.com/hiank/think/net"
	"github.com/hiank/thinkend/war"
)


func main() {

	defer glog.Infoln("close war-server")

	flag.Parse()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	ctx := context.Background()
	go serveK8sLink(ctx, wg)

	wg.Wait()
}


func serveK8sLink(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()

	h := war.NewHandler(ctx)
	defer h.Close()		//NOTE: 关闭Handler的context
	glog.Infoln("before servek8s")
	if e := net.ServeK8s(ctx, "", h); e != nil {

		glog.Fatalln("serve k8s link error : ", e)
	}
	glog.Infoln("after servek8s")
}