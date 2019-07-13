package main

import (
	"sync"
	"flag"
	"strconv"
	"net/url"
	"net/http"
	"github.com/gorilla/websocket"
	"github.com/hiank/think/pb"
	"github.com/golang/protobuf/ptypes"
	"fmt"
	war_pb "github.com/hiank/thinkend/war/proto"
)

func main() {

	var num int
	flag.IntVar(&num, "n", 4, "num of client")
	flag.Parse()

	wait := new(sync.WaitGroup)
	wait.Add(num)
	max := num + 1001
	for i:=1001; i < max; i++ {

		go MakeTank(wait, i)
	} 
	wait.Wait()
}


//MakeTank 创建一个客户端
func MakeTank(wait *sync.WaitGroup, tokenNum int) {

	defer wait.Done()

	token := strconv.Itoa(tokenNum)
	conn, _, err := dail(token)
	if err != nil {
		fmt.Print(err)
		return
	}
	defer conn.Close()

	//NOTE: do war request
	anyMsg, err := ptypes.MarshalAny(&war_pb.S_War_Want{Type: war_pb.War_Type_TEAMWINNER})
	if err != nil {
		fmt.Printf("marshal any error : %v\n", err)
		return
	}

	buf, err := pb.AnyEncode(anyMsg)
	if err != nil {
		fmt.Printf("encode error : %v\n", err)
		return
	}

	if err := conn.WriteMessage(websocket.BinaryMessage, buf); err != nil {

		fmt.Printf("send want message error : %v\n", err)
		return
	}
	
	for {

		fmt.Println("before read message")
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Printf("read error : %v\n", err)
			return
		}
		fmt.Printf("received: %s\n", message)
	}
}


func dail(token string) (*websocket.Conn, *http.Response, error) {

	addr := "192.168.137.222:30250"
	fmt.Printf("address : %s\n", addr)

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	return websocket.DefaultDialer.Dial(u.String(), http.Header{"token": {token}})
}




