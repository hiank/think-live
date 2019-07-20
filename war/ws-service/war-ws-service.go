package main

import (
	"github.com/hiank/think/pool"
	"github.com/hiank/thinkend/war"
	"github.com/hiank/think/net"
	"github.com/hiank/think/net/ws"
	"context"
	"sync"
	"flag"
	"github.com/golang/glog"
)

func main() {

	defer glog.Infoln("close war-ws-server")

	flag.Parse()

	wg := new(sync.WaitGroup)
	wg.Add(1)

	ctx := context.Background()
	net.Init(ctx)
	go serveWs(ctx, wg)

	wg.Wait()
}


func serveWs(ctx context.Context, wg *sync.WaitGroup) {

	defer wg.Done()

	war.SetPoolGetter(war.PoolGetter(func() *pool.Pool {

		return ws.GetWSPool()
	}))
	// war.SetNetPool(ws.GetWSPool())
	h := war.NewWSHandler(ctx)
	defer h.Close()		//NOTE: 关闭Handler的context
	
	glog.Infoln("before serveWs")	
	if e := ws.ListenAndServeWS(ctx, "", h); e != nil {

		glog.Fatalln("serve ws error : ", e)
	}
	glog.Infoln("after serveWs")
}
