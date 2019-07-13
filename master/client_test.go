package master_test

import (
	"net/http"
	"github.com/gorilla/websocket"
	"net/url"
	"fmt"
)

func Dial() (*websocket.Conn, *http.Response, error) {

	addr := "192.168.137.222:30250"
	fmt.Println("address : " + addr)

	u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
	return websocket.DefaultDialer.Dial(u.String(), nil)
}